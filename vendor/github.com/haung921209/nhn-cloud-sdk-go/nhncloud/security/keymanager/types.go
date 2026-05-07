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
	KeyStoreID  string `json:"keyStoreId"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	CreatedAt   string `json:"createdAt,omitempty"`
	UpdatedAt   string `json:"updatedAt,omitempty"`
	KeyCount    int    `json:"keyCount,omitempty"`
	Status      string `json:"status,omitempty"`
}

type ListKeyStoresOutput struct {
	APIResponse
	Body struct {
		KeyStores []KeyStore `json:"keystores"`
	} `json:"body"`
}

type GetKeyStoreOutput struct {
	APIResponse
	Body struct {
		KeyStore KeyStore `json:"keystore"`
	} `json:"body"`
}

// ============== Key Types ==============

type Key struct {
	KeyID          string `json:"keyId"`
	KeyStoreID     string `json:"keyStoreId,omitempty"`
	Name           string `json:"name"`
	Description    string `json:"description,omitempty"`
	KeyType        string `json:"keyType"` // secrets, symmetric-keys, asymmetric-keys
	KeyAlgorithm   string `json:"keyAlgorithm,omitempty"`
	KeySize        int    `json:"keySize,omitempty"`
	CreatedAt      string `json:"createdAt,omitempty"`
	UpdatedAt      string `json:"updatedAt,omitempty"`
	ExpirationDate string `json:"expirationDate,omitempty"`
	DeletionDate   string `json:"deletionDate,omitempty"`
	Status         string `json:"status,omitempty"`
	RotationPeriod int    `json:"rotationPeriod,omitempty"`
}

type ListKeysOutput struct {
	APIResponse
	Body struct {
		Keys []Key `json:"keys"`
	} `json:"body"`
}

type GetKeyOutput struct {
	APIResponse
	Body struct {
		Key Key `json:"key"`
	} `json:"body"`
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
		KeyValue string `json:"keyValue"`
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
		PlainDataKey     string `json:"plainDataKey"`
		EncryptedDataKey string `json:"encryptedDataKey"`
	} `json:"body"`
}

// ============== Asymmetric Key Types ==============

type GetPrivateKeyOutput struct {
	APIResponse
	Body struct {
		PrivateKey string `json:"privateKey"`
	} `json:"body"`
}

type GetPublicKeyOutput struct {
	APIResponse
	Body struct {
		PublicKey string `json:"publicKey"`
	} `json:"body"`
}

type SignInput struct {
	Data string `json:"data"` // Base64 encoded data to sign
}

type SignOutput struct {
	APIResponse
	Body struct {
		Signature string `json:"signature"` // Base64 encoded signature
	} `json:"body"`
}

type VerifyInput struct {
	Data      string `json:"data"`      // Base64 encoded original data
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
	Algorithm      string `json:"algorithm,omitempty"`      // RSA2048, RSA4096, EC_P256, EC_P384 for asymmetric
	Secret         string `json:"secret,omitempty"`         // For secrets
	RotationPeriod int    `json:"rotationPeriod,omitempty"` // days
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
		KeyID        string `json:"keyId"`
		DeletionDate string `json:"deletionDate,omitempty"`
	} `json:"body"`
}

// ============== Client Info Types ==============

type GetClientInfoOutput struct {
	APIResponse
	Body struct {
		AppKey     string `json:"appKey"`
		IPAddress  string `json:"ipAddress,omitempty"`
		MACAddress string `json:"macAddress,omitempty"`
	} `json:"body"`
}
