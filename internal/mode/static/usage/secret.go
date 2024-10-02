package usage

import (
	"sync"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

// Secret implements the SecretStorer interface.
type Secret struct {
	secret *v1.Secret
	lock   *sync.Mutex
	nsName types.NamespacedName
}

// NewUsageSecret creates a new Secret wrapper.
func NewUsageSecret(nsName types.NamespacedName) *Secret {
	return &Secret{
		lock:   &sync.Mutex{},
		nsName: nsName,
	}
}

// Set stores the updated Secre.
func (s *Secret) Set(secret *v1.Secret) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.secret = secret
}

// Delete nullifies the Secret value.
func (s *Secret) Delete() {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.secret = nil
}

// GetNSName returns the namespaced name of the Secret.
func (s *Secret) GetNSName() types.NamespacedName {
	return s.nsName
}

// GetJWT returns the base64 encoded JWT from the Secret.
func (s *Secret) GetJWT() []byte {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.secret != nil {
		return s.secret.Data["license.jwt"]
	}

	return nil
}
