package keymanager

// APIResponse wrapper for Key Manager API responses
type APIResponse struct {
	Header struct {
		ResultCode    int    `json:"resultCode"`
		ResultMessage string `json:"resultMessage"`
		IsSuccessful  bool   `json:"isSuccessful"`
	} `json:"header"`
}

// ============== Key Store Types ==============

type KeyStore struct {
	KeyStoreID         int64  `json:"keyStoreId"`
	Name               string `json:"name"`
	Description        string `json:"description,omitempty"`
	IP4AuthUse         string `json:"ip4AuthUse,omitempty"` // Y or N
	MACAuthUse         string `json:"macAuthUse,omitempty"`
	CertificateAuthUse string `json:"certificateAuthUse,omitempty"`
	CreationUser       string `json:"creationUser,omitempty"`
	CreationDatetime   string `json:"creationDatetime,omitempty"`
	LastChangeUser     string `json:"lastChangeUser,omitempty"`
	LastChangeDatetime string `json:"lastChangeDatetime,omitempty"`
}

type ListKeyStoresOutput struct {
	APIResponse
	Body struct {
		KeyStores []KeyStore `json:"keyStoreList"`
	} `json:"body"`
}

type GetKeyStoreOutput struct {
	APIResponse
	Body KeyStore `json:"body"`
}

// ============== Key Types ==============

type Key struct {
	KeyID                  string `json:"keyId"`
	Name                   string `json:"name"`
	Description            string `json:"description,omitempty"`
	KeyType                string `json:"keyType"` // SECRET, SYMMETRIC_KEY, ASYMMETRIC_KEY
	KeyStatus              string `json:"keyStatus,omitempty"`
	CurrentKeyValueVersion int    `json:"currentKeyValueVersion,omitempty"`
	AutoRotationPeriod     int    `json:"autoRotationPeriod,omitempty"`
	NextAutoRotationDate   string `json:"nextAutoRotationDate,omitempty"`
	LastAccessDatetime     string `json:"lastAccessDatetime,omitempty"`
	DeletionDatetime       string `json:"deletionDatetime,omitempty"`
	CreationUser           string `json:"creationUser,omitempty"`
	CreationDatetime       string `json:"creationDatetime,omitempty"`
	LastChangeUser         string `json:"lastChangeUser,omitempty"`
	LastChangeDatetime     string `json:"lastChangeDatetime,omitempty"`
}

type ListKeysOutput struct {
	APIResponse
	Body struct {
		Keys []Key `json:"keyList"`
	} `json:"body"`
}

type GetKeyOutput struct {
	APIResponse
	Body Key `json:"body"`
}

// ============== Secret Types ==============

type GetSecretOutput struct {
	APIResponse
	Body struct {
		Secret string `json:"secret"`
	} `json:"body"`
}

// ============== Symmetric Key Types ==============

type GetSymmetricKeyOutput struct {
	APIResponse
	Body struct {
		SymmetricKey string `json:"symmetricKey"`
	} `json:"body"`
}

type EncryptInput struct {
	Plaintext string `json:"plaintext"`
	AAD       string `json:"aad,omitempty"` // Additional Authenticated Data
}

type EncryptOutput struct {
	APIResponse
	Body struct {
		Ciphertext string `json:"ciphertext"`
		IV         string `json:"iv,omitempty"`
		Tag        string `json:"tag,omitempty"`
	} `json:"body"`
}

type DecryptInput struct {
	Ciphertext string `json:"ciphertext"`
	IV         string `json:"iv,omitempty"`
	Tag        string `json:"tag,omitempty"`
	AAD        string `json:"aad,omitempty"`
}

type DecryptOutput struct {
	APIResponse
	Body struct {
		Plaintext string `json:"plaintext"`
	} `json:"body"`
}

type CreateLocalKeyOutput struct {
	APIResponse
	Body struct {
		PlainDataKey     string `json:"localKeyPlaintext"`
		EncryptedDataKey string `json:"localKeyCiphertext"`
	} `json:"body"`
}

// ============== Asymmetric Key Types ==============

// AsymmetricKey is the body of the private and public key reads: Key is the
// key as a byte list, EncodedKey the same key in Base64 (DER).
type AsymmetricKey struct {
	KeyType    string `json:"keyType"` // PrivateKey or PublicKey
	Key        string `json:"key"`
	EncodedKey string `json:"encodedKey"`
	KeyVersion int    `json:"keyVersion"`
}

type GetPrivateKeyOutput struct {
	APIResponse
	Body AsymmetricKey `json:"body"`
}

type GetPublicKeyOutput struct {
	APIResponse
	Body AsymmetricKey `json:"body"`
}

type SignInput struct {
	Data string `json:"plaintext"` // Base64 encoded data to sign
}

type SignOutput struct {
	APIResponse
	Body struct {
		Signature string `json:"signature"` // Base64 encoded signature
	} `json:"body"`
}

type VerifyInput struct {
	Data      string `json:"plaintext"` // Base64 encoded original data
	Signature string `json:"signature"` // Base64 encoded signature to verify
}

type VerifyOutput struct {
	APIResponse
	Body struct {
		Result bool `json:"result"` // true if signature is valid
	} `json:"body"`
}

// ============== Key Management Types ==============

type CreateKeyInput struct {
	KeyStoreName   string `json:"keyStoreName"`
	Name           string `json:"name"`
	Description    string `json:"description,omitempty"`
	KeyAlgorithm   string `json:"keyAlgorithm,omitempty"` // AES256 for symmetric
	KeySize        int    `json:"keySize,omitempty"`
	Algorithm      string `json:"algorithm,omitempty"`          // RSA2048, RSA4096, EC_P256, EC_P384 for asymmetric
	Secret         string `json:"secretValue,omitempty"`        // For secrets
	RotationPeriod int    `json:"autoRotationPeriod,omitempty"` // days
}

type CreateKeyOutput struct {
	APIResponse
	Body struct {
		KeyID string `json:"keyId"`
		Key   Key    `json:"key,omitempty"`
	} `json:"body"`
}

type DeleteKeyInput struct {
	RequestDeletionDate string `json:"requestDeletionDate,omitempty"` // ISO 8601 format
}

type DeleteKeyOutput struct {
	APIResponse
	Body struct {
		KeyID            string `json:"keyId"`
		DeletionDatetime string `json:"deletionDateTime,omitempty"`
	} `json:"body"`
}

// ============== Client Info Types ==============

type GetClientInfoOutput struct {
	APIResponse
	Body struct {
		ClientIP              string `json:"clientIp"`
		ClientMACHeader       string `json:"clientMacHeader"`
		ClientSentCertificate bool   `json:"clientSentCertificate"`
	} `json:"body"`
}
