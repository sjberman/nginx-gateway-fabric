package usage

import (
	"testing"

	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func TestSet(t *testing.T) {
	t.Parallel()
	store := NewUsageSecret(types.NamespacedName{})
	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "secret",
			Namespace: "custom",
		},
	}

	g := NewWithT(t)
	g.Expect(store.secret).To(BeNil())

	store.Set(secret)
	g.Expect(store.secret).To(Equal(secret))
}

func TestDelete(t *testing.T) {
	t.Parallel()
	store := NewUsageSecret(types.NamespacedName{})
	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "secret",
			Namespace: "custom",
		},
	}

	g := NewWithT(t)
	store.Set(secret)
	g.Expect(store.secret).To(Equal(secret))

	store.Delete()
	g.Expect(store.secret).To(BeNil())
}

func TestGetNSName(t *testing.T) {
	t.Parallel()

	nsName := types.NamespacedName{Name: "secret", Namespace: "custom"}
	store := NewUsageSecret(nsName)

	g := NewWithT(t)
	g.Expect(store.GetNSName()).To(Equal(nsName))
}

func TestGetJWT(t *testing.T) {
	t.Parallel()
	store := NewUsageSecret(types.NamespacedName{})
	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "secret",
			Namespace: "custom",
		},
		Data: map[string][]byte{
			"license.jwt": []byte("license"),
		},
	}

	g := NewWithT(t)

	jwt := store.GetJWT()
	g.Expect(jwt).To(BeNil())

	store.Set(secret)

	jwt = store.GetJWT()
	g.Expect(jwt).To(Equal([]byte("license")))
}
