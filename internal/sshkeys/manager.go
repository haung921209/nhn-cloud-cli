package sshkeys

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

const (
	sshKeysDir   = ".nhncloud/ssh-keys"
	metadataFile = "keys.json"
)

type KeyInfo struct {
	Name        string    `json:"name"`
	Path        string    `json:"path"`
	Fingerprint string    `json:"fingerprint"`
	Type        string    `json:"type"`
	PublicKey   string    `json:"public_key,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	LastUsed    time.Time `json:"last_used,omitempty"`
}

type Manager struct {
	baseDir string
}

func NewManager() *Manager {
	home, _ := os.UserHomeDir()
	return &Manager{
		baseDir: filepath.Join(home, sshKeysDir),
	}
}

func (m *Manager) ensureDir() error {
	return os.MkdirAll(m.baseDir, 0700)
}

func (m *Manager) List() ([]KeyInfo, error) {
	if err := m.ensureDir(); err != nil {
		return nil, err
	}

	metaPath := filepath.Join(m.baseDir, metadataFile)
	data, err := os.ReadFile(metaPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []KeyInfo{}, nil
		}
		return nil, err
	}

	var keys []KeyInfo
	if err := json.Unmarshal(data, &keys); err != nil {
		return nil, err
	}

	return keys, nil
}

func (m *Manager) saveMetadata(keys []KeyInfo) error {
	metaPath := filepath.Join(m.baseDir, metadataFile)
	data, err := json.MarshalIndent(keys, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(metaPath, data, 0600)
}

func (m *Manager) Import(name, filePath string) (*KeyInfo, error) {
	if err := m.ensureDir(); err != nil {
		return nil, err
	}

	keyData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read key file: %w", err)
	}

	keyInfo, err := m.parsePrivateKey(name, keyData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	destPath := filepath.Join(m.baseDir, name+".pem")
	if err := os.WriteFile(destPath, keyData, 0600); err != nil {
		return nil, fmt.Errorf("failed to save key: %w", err)
	}

	keyInfo.Path = destPath
	keyInfo.CreatedAt = time.Now()

	keys, _ := m.List()

	var newKeys []KeyInfo
	for _, k := range keys {
		if k.Name != name {
			newKeys = append(newKeys, k)
		}
	}
	newKeys = append(newKeys, *keyInfo)

	if err := m.saveMetadata(newKeys); err != nil {
		os.Remove(destPath)
		return nil, err
	}

	return keyInfo, nil
}

func (m *Manager) parsePrivateKey(name string, keyData []byte) (*KeyInfo, error) {
	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block")
	}

	var keyType string
	var publicKeyStr string
	var fingerprint string

	switch block.Type {
	case "RSA PRIVATE KEY":
		keyType = "RSA"
		key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		publicKey, err := ssh.NewPublicKey(&key.PublicKey)
		if err != nil {
			return nil, err
		}
		publicKeyStr = string(ssh.MarshalAuthorizedKey(publicKey))
		fingerprint = ssh.FingerprintLegacyMD5(publicKey)

	case "PRIVATE KEY":
		keyType = "PKCS8"
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			signer, err := ssh.ParsePrivateKey(keyData)
			if err != nil {
				return nil, fmt.Errorf("unsupported private key format: %w", err)
			}
			keyType = "OPENSSH"
			publicKeyStr = string(ssh.MarshalAuthorizedKey(signer.PublicKey()))
			fingerprint = ssh.FingerprintLegacyMD5(signer.PublicKey())
		} else {
			if rsaKey, ok := key.(*rsa.PrivateKey); ok {
				keyType = "RSA (PKCS8)"
				publicKey, err := ssh.NewPublicKey(&rsaKey.PublicKey)
				if err != nil {
					return nil, err
				}
				publicKeyStr = string(ssh.MarshalAuthorizedKey(publicKey))
				fingerprint = ssh.FingerprintLegacyMD5(publicKey)
			}
		}

	case "OPENSSH PRIVATE KEY":
		keyType = "OPENSSH"
		signer, err := ssh.ParsePrivateKey(keyData)
		if err != nil {
			return nil, err
		}
		publicKeyStr = string(ssh.MarshalAuthorizedKey(signer.PublicKey()))
		fingerprint = ssh.FingerprintLegacyMD5(signer.PublicKey())

	default:
		return nil, fmt.Errorf("unsupported key type: %s", block.Type)
	}

	return &KeyInfo{
		Name:        name,
		Type:        keyType,
		PublicKey:   strings.TrimSpace(publicKeyStr),
		Fingerprint: fingerprint,
	}, nil
}

func (m *Manager) Get(name string) (*KeyInfo, error) {
	keys, err := m.List()
	if err != nil {
		return nil, err
	}

	for _, key := range keys {
		if key.Name == name {
			return &key, nil
		}
	}

	return nil, fmt.Errorf("SSH key '%s' not found", name)
}

func (m *Manager) Remove(name string) error {
	keys, err := m.List()
	if err != nil {
		return err
	}

	var newKeys []KeyInfo
	var found bool
	for _, key := range keys {
		if key.Name == name {
			found = true
			if err := os.Remove(key.Path); err != nil && !os.IsNotExist(err) {
				return err
			}
		} else {
			newKeys = append(newKeys, key)
		}
	}

	if !found {
		return fmt.Errorf("SSH key '%s' not found", name)
	}

	return m.saveMetadata(newKeys)
}

func (m *Manager) Connect(keyName, target string) error {
	keyInfo, err := m.Get(keyName)
	if err != nil {
		return err
	}

	keys, _ := m.List()
	for i := range keys {
		if keys[i].Name == keyName {
			keys[i].LastUsed = time.Now()
			m.saveMetadata(keys)
			break
		}
	}

	cmd := exec.Command("ssh", "-i", keyInfo.Path, "-o", "StrictHostKeyChecking=no", target)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (m *Manager) Export(keyName, destPath string) error {
	keyInfo, err := m.Get(keyName)
	if err != nil {
		return err
	}

	keyData, err := os.ReadFile(keyInfo.Path)
	if err != nil {
		return fmt.Errorf("failed to read key: %w", err)
	}

	if err := os.WriteFile(destPath, keyData, 0600); err != nil {
		return fmt.Errorf("failed to write key: %w", err)
	}

	return nil
}
