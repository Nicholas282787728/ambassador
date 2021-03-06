package snapshot

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	amb "github.com/datawire/ambassador/pkg/api/getambassador.io/v2"
	"github.com/datawire/ambassador/pkg/kates"
	"github.com/datawire/ambassador/pkg/watt"
)

const ApiVersion = "v1"
const ContentTypeJSON = "application/json"

// SecretRef is a secret reference -- basically, a namespace/name pair.
type SecretRef struct {
	Namespace string
	Name      string
}

// The snapshot type represents a complete configuration snapshot as sent to
// diagd.
type Snapshot struct {
	// The Kubernetes field contains all the ambassador inputs from kubernetes.
	Kubernetes *KubernetesSnapshot
	// The Consul field contains endpoint data for any mappings setup to use a
	// consul resolver.
	Consul *watt.ConsulSnapshot
	// The Deltas field contains a list of deltas to indicate what has changed
	// since the prior snapshot. This is only computed for the Kubernetes
	// portion of the snapshot. Changes in the Consul endpoint data are not
	// reflected in this field.
	Deltas []*kates.Delta
	// The Invalid field contains any kubernetes resources that have failed
	// validation.
	Invalid []*kates.Unstructured
	Raw     json.RawMessage `json:"-"`
}

type KubernetesSnapshot struct {
	// k8s resources
	IngressClasses []*kates.IngressClass `json:"ingressclasses"`
	Ingresses      []*kates.Ingress      `json:"ingresses"`
	Services       []*kates.Service      `json:"service"`

	K8sEndpoints []*kates.Endpoints `json:"-"`         // Endpoints from Kubernetes
	Endpoints    []*kates.Endpoints `json:"Endpoints"` // Endpoints we'll actually feed to Ambassador

	// ambassador resources
	Hosts       []*amb.Host       `json:"Host"`
	Mappings    []*amb.Mapping    `json:"Mapping"`
	TCPMappings []*amb.TCPMapping `json:"TCPMapping"`
	Modules     []*amb.Module     `json:"Module"`
	TLSContexts []*amb.TLSContext `json:"TLSContext"`

	// plugin services
	AuthServices      []*amb.AuthService      `json:"AuthService"`
	RateLimitServices []*amb.RateLimitService `json:"RateLimitService"`
	LogServices       []*amb.LogService       `json:"LogService"`
	TracingServices   []*amb.TracingService   `json:"TracingService"`
	DevPortals        []*amb.DevPortal        `json:"DevPortal"`

	// resolvers
	ConsulResolvers             []*amb.ConsulResolver             `json:"ConsulResolver"`
	KubernetesEndpointResolvers []*amb.KubernetesEndpointResolver `json:"KubernetesEndpointResolver"`
	KubernetesServiceResolvers  []*amb.KubernetesServiceResolver  `json:"KubernetesServiceResolver"`

	// It is safe to ignore AmbassadorInstallation, ambassador doesn't need to look at those, just
	// the operator.

	KNativeClusterIngresses []*kates.Unstructured `json:"clusteringresses.networking.internal.knative.dev,omitempty"`
	KNativeIngresses        []*kates.Unstructured `json:"ingresses.networking.internal.knative.dev,omitempty"`

	K8sSecrets []*kates.Secret             `json:"-"`      // Secrets from Kubernetes
	FSSecrets  map[SecretRef]*kates.Secret `json:"-"`      // Secrets from the filesystem
	Secrets    []*kates.Secret             `json:"secret"` // Secrets we'll feed to Ambassador

	Annotations []kates.Object `json:"-"`
}

func (a *KubernetesSnapshot) Render() string {
	result := &strings.Builder{}
	v := reflect.ValueOf(a)
	t := v.Type().Elem()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		result.WriteString(fmt.Sprintf("%s: %d\n", f.Name, reflect.Indirect(v).Field(i).Len()))
	}
	return result.String()
}
