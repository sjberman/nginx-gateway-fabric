package graph

import (
	"testing"

	. "github.com/onsi/gomega"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	"sigs.k8s.io/gateway-api/apis/v1alpha2"

	ngfAPI "github.com/nginx/nginx-gateway-fabric/v2/apis/v1alpha2"
	"github.com/nginx/nginx-gateway-fabric/v2/internal/controller/state/conditions"
	"github.com/nginx/nginx-gateway-fabric/v2/internal/framework/helpers"
)

func createTLSRoute(
	hostname gatewayv1.Hostname,
	rules []v1alpha2.TLSRouteRule,
	parentRefs []gatewayv1.ParentReference,
) *v1alpha2.TLSRoute {
	return &v1alpha2.TLSRoute{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test",
			Name:      "tr",
		},
		Spec: v1alpha2.TLSRouteSpec{
			CommonRouteSpec: gatewayv1.CommonRouteSpec{
				ParentRefs: parentRefs,
			},
			Hostnames: []gatewayv1.Hostname{hostname},
			Rules:     rules,
		},
	}
}

func TestBuildTLSRoute(t *testing.T) {
	t.Parallel()

	parentRef := gatewayv1.ParentReference{
		Namespace:   helpers.GetPointer[gatewayv1.Namespace]("test"),
		Name:        "gateway",
		SectionName: helpers.GetPointer[gatewayv1.SectionName]("l1"),
	}

	createGateway := func() *Gateway {
		return &Gateway{
			Source: &gatewayv1.Gateway{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "test",
					Name:      "gateway",
				},
			},
			Valid: true,
		}
	}

	modGateway := func(gw *Gateway, mod func(*Gateway) *Gateway) *Gateway {
		return mod(gw)
	}

	parentRefGraph := ParentRef{
		SectionName: helpers.GetPointer[gatewayv1.SectionName]("l1"),
		Gateway: &ParentRefGateway{
			NamespacedName: types.NamespacedName{
				Namespace: "test",
				Name:      "gateway",
			},
		},
	}
	duplicateParentRefsGtr := createTLSRoute(
		"hi.example.com",
		nil,
		[]gatewayv1.ParentReference{
			parentRef,
			parentRef,
		},
	)
	noParentRefsGtr := createTLSRoute(
		"hi.example.com",
		nil,
		[]gatewayv1.ParentReference{},
	)
	invalidHostnameGtr := createTLSRoute(
		"hi....com",
		nil,
		[]gatewayv1.ParentReference{
			parentRef,
		},
	)
	noRulesGtr := createTLSRoute(
		"app.example.com",
		nil,
		[]gatewayv1.ParentReference{
			parentRef,
		},
	)
	backedRefDNEGtr := createTLSRoute(
		"app.example.com",
		[]v1alpha2.TLSRouteRule{
			{
				BackendRefs: []gatewayv1.BackendRef{
					{
						BackendObjectReference: gatewayv1.BackendObjectReference{
							Name: "hi",
							Port: helpers.GetPointer[gatewayv1.PortNumber](80),
						},
					},
				},
			},
		},
		[]gatewayv1.ParentReference{
			parentRef,
		},
	)

	wrongBackendRefGroupGtr := createTLSRoute(
		"app.example.com",
		[]v1alpha2.TLSRouteRule{
			{
				BackendRefs: []gatewayv1.BackendRef{
					{
						BackendObjectReference: gatewayv1.BackendObjectReference{
							Name:  "hi",
							Port:  helpers.GetPointer[gatewayv1.PortNumber](80),
							Group: helpers.GetPointer[gatewayv1.Group]("wrong"),
						},
					},
				},
			},
		},
		[]gatewayv1.ParentReference{
			parentRef,
		},
	)

	wrongBackendRefKindGtr := createTLSRoute(
		"app.example.com",
		[]v1alpha2.TLSRouteRule{
			{
				BackendRefs: []gatewayv1.BackendRef{
					{
						BackendObjectReference: gatewayv1.BackendObjectReference{
							Name: "hi",
							Port: helpers.GetPointer[gatewayv1.PortNumber](80),
							Kind: helpers.GetPointer[gatewayv1.Kind]("not service"),
						},
					},
				},
			},
		},
		[]gatewayv1.ParentReference{
			parentRef,
		},
	)

	diffNsBackendRef := createTLSRoute("app.example.com",
		[]v1alpha2.TLSRouteRule{
			{
				BackendRefs: []gatewayv1.BackendRef{
					{
						BackendObjectReference: gatewayv1.BackendObjectReference{
							Name:      "hi",
							Port:      helpers.GetPointer[gatewayv1.PortNumber](80),
							Namespace: helpers.GetPointer[gatewayv1.Namespace]("diff"),
						},
					},
				},
			},
		},
		[]gatewayv1.ParentReference{
			parentRef,
		},
	)

	portNilBackendRefGtr := createTLSRoute("app.example.com",
		[]v1alpha2.TLSRouteRule{
			{
				BackendRefs: []gatewayv1.BackendRef{
					{
						BackendObjectReference: gatewayv1.BackendObjectReference{
							Name: "hi",
						},
					},
				},
			},
		},
		[]gatewayv1.ParentReference{
			parentRef,
		},
	)

	ipFamilyMismatchGtr := createTLSRoute(
		"app.example.com",
		[]v1alpha2.TLSRouteRule{
			{
				BackendRefs: []gatewayv1.BackendRef{
					{
						BackendObjectReference: gatewayv1.BackendObjectReference{
							Name: "hi",
							Port: helpers.GetPointer[gatewayv1.PortNumber](80),
						},
					},
				},
			},
		},
		[]gatewayv1.ParentReference{
			parentRef,
		},
	)

	validRefSameNs := createTLSRoute("app.example.com",
		[]v1alpha2.TLSRouteRule{
			{
				BackendRefs: []gatewayv1.BackendRef{
					{
						BackendObjectReference: gatewayv1.BackendObjectReference{
							Name:      "hi",
							Port:      helpers.GetPointer[gatewayv1.PortNumber](80),
							Namespace: helpers.GetPointer[gatewayv1.Namespace]("test"),
						},
					},
				},
			},
		},
		[]gatewayv1.ParentReference{
			parentRef,
		},
	)

	svcNsName := types.NamespacedName{
		Namespace: "test",
		Name:      "hi",
	}

	diffSvcNsName := types.NamespacedName{
		Namespace: "diff",
		Name:      "hi",
	}

	createSvc := func(name string, port int32) *apiv1.Service {
		return &apiv1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "test",
				Name:      name,
			},
			Spec: apiv1.ServiceSpec{
				Ports: []apiv1.ServicePort{
					{Port: port},
				},
			},
		}
	}

	createSvcWithAppProtocol := func(name, appProtocol string, port int32) *apiv1.Service {
		return &apiv1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "test",
				Name:      name,
			},
			Spec: apiv1.ServiceSpec{
				Ports: []apiv1.ServicePort{
					{Port: port, AppProtocol: &appProtocol},
				},
			},
		}
	}

	diffNsSvc := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "diff",
			Name:      "hi",
		},
		Spec: apiv1.ServiceSpec{
			Ports: []apiv1.ServicePort{
				{Port: 80},
			},
		},
	}

	ipv4Svc := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test",
			Name:      "hi",
		},
		Spec: apiv1.ServiceSpec{
			IPFamilies: []apiv1.IPFamily{
				apiv1.IPv4Protocol,
			},
			Ports: []apiv1.ServicePort{
				{Port: 80},
			},
		},
	}

	alwaysTrueRefGrantResolver := func(_ toResource) bool { return true }
	alwaysFalseRefGrantResolver := func(_ toResource) bool { return false }

	tests := []struct {
		expected *L4Route
		gtr      *v1alpha2.TLSRoute
		services map[types.NamespacedName]*apiv1.Service
		resolver func(resource toResource) bool
		gateway  *Gateway
		name     string
	}{
		{
			gtr: duplicateParentRefsGtr,
			expected: &L4Route{
				Source: duplicateParentRefsGtr,
				Valid:  false,
			},
			gateway:  createGateway(),
			services: map[types.NamespacedName]*apiv1.Service{},
			resolver: alwaysTrueRefGrantResolver,
			name:     "duplicate parent refs",
		},
		{
			gtr:      noParentRefsGtr,
			expected: nil,
			gateway:  createGateway(),
			services: map[types.NamespacedName]*apiv1.Service{},
			resolver: alwaysTrueRefGrantResolver,
			name:     "no parent refs",
		},
		{
			gtr: invalidHostnameGtr,
			expected: &L4Route{
				Source:     invalidHostnameGtr,
				ParentRefs: []ParentRef{parentRefGraph},
				Conditions: []conditions.Condition{conditions.NewRouteUnsupportedValue(
					"spec.hostnames[0]: Invalid value: \"hi....com\": a lowercase RFC 1" +
						"123 subdomain must consist of lower case alphanumeric characters" +
						", '-' or '.', and must start and end with an alphanumeric charac" +
						"ter (e.g. 'example.com', regex used for validation is '[a-z0-9](" +
						"[-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*')",
				)},
				Valid: false,
			},
			gateway:  createGateway(),
			services: map[types.NamespacedName]*apiv1.Service{},
			resolver: alwaysTrueRefGrantResolver,
			name:     "invalid hostname",
		},
		{
			gtr: noRulesGtr,
			expected: &L4Route{
				Source:     noRulesGtr,
				ParentRefs: []ParentRef{parentRefGraph},
				Spec: L4RouteSpec{
					Hostnames: []gatewayv1.Hostname{
						"app.example.com",
					},
				},
				Conditions: []conditions.Condition{conditions.NewRouteBackendRefUnsupportedValue(
					"Must have exactly one Rule and BackendRef",
				)},
				Valid: false,
			},
			gateway:  createGateway(),
			services: map[types.NamespacedName]*apiv1.Service{},
			resolver: alwaysTrueRefGrantResolver,
			name:     "invalid rule",
		},
		{
			gtr: validRefSameNs,
			expected: &L4Route{
				Source:     validRefSameNs,
				ParentRefs: []ParentRef{parentRefGraph},
				Spec: L4RouteSpec{
					Hostnames: []gatewayv1.Hostname{
						"app.example.com",
					},
					BackendRef: BackendRef{
						SvcNsName:          svcNsName,
						ServicePort:        apiv1.ServicePort{Port: 80, AppProtocol: helpers.GetPointer(AppProtocolTypeH2C)},
						Valid:              false,
						InvalidForGateways: map[types.NamespacedName]conditions.Condition{},
					},
				},
				Attachable: true,
				Valid:      true,
				Conditions: []conditions.Condition{conditions.NewRouteBackendRefUnsupportedProtocol(
					"route type tls does not support service port appProtocol kubernetes.io/h2c",
				)},
			},
			gateway: createGateway(),
			services: map[types.NamespacedName]*apiv1.Service{
				svcNsName: createSvcWithAppProtocol("hi", AppProtocolTypeH2C, 80),
			},
			resolver: alwaysTrueRefGrantResolver,
			name:     "invalid service port appProtocol h2c",
		},
		{
			gtr: validRefSameNs,
			expected: &L4Route{
				Source:     validRefSameNs,
				ParentRefs: []ParentRef{parentRefGraph},
				Spec: L4RouteSpec{
					Hostnames: []gatewayv1.Hostname{
						"app.example.com",
					},
					BackendRef: BackendRef{
						SvcNsName:          svcNsName,
						ServicePort:        apiv1.ServicePort{Port: 80, AppProtocol: helpers.GetPointer(AppProtocolTypeWS)},
						Valid:              false,
						InvalidForGateways: map[types.NamespacedName]conditions.Condition{},
					},
				},
				Attachable: true,
				Valid:      true,
				Conditions: []conditions.Condition{conditions.NewRouteBackendRefUnsupportedProtocol(
					"route type tls does not support service port appProtocol kubernetes.io/ws",
				)},
			},
			gateway: createGateway(),
			services: map[types.NamespacedName]*apiv1.Service{
				svcNsName: createSvcWithAppProtocol("hi", AppProtocolTypeWS, 80),
			},
			resolver: alwaysTrueRefGrantResolver,
			name:     "invalid service port appProtocol WS",
		},
		{
			gtr: backedRefDNEGtr,
			expected: &L4Route{
				Source:     backedRefDNEGtr,
				ParentRefs: []ParentRef{parentRefGraph},
				Spec: L4RouteSpec{
					Hostnames: []gatewayv1.Hostname{
						"app.example.com",
					},
					BackendRef: BackendRef{
						SvcNsName: types.NamespacedName{
							Namespace: "test",
							Name:      "hi",
						},
						Valid:              false,
						InvalidForGateways: map[types.NamespacedName]conditions.Condition{},
					},
				},
				Conditions: []conditions.Condition{conditions.NewRouteBackendRefRefBackendNotFound(
					"spec.rules[0].backendRefs[0].name: Not found: \"hi\"",
				)},
				Attachable: true,
				Valid:      true,
			},
			gateway:  createGateway(),
			services: map[types.NamespacedName]*apiv1.Service{},
			resolver: alwaysTrueRefGrantResolver,
			name:     "BackendRef not found",
		},
		{
			gtr: wrongBackendRefGroupGtr,
			expected: &L4Route{
				Source:     wrongBackendRefGroupGtr,
				ParentRefs: []ParentRef{parentRefGraph},
				Spec: L4RouteSpec{
					Hostnames: []gatewayv1.Hostname{
						"app.example.com",
					},
					BackendRef: BackendRef{
						Valid:              false,
						InvalidForGateways: map[types.NamespacedName]conditions.Condition{},
					},
				},
				Conditions: []conditions.Condition{conditions.NewRouteBackendRefInvalidKind(
					"spec.rules[0].backendRefs[0].group:" +
						" Unsupported value: \"wrong\": supported values: \"core\", \"\"",
				)},
				Attachable: true,
				Valid:      true,
			},
			gateway: createGateway(),
			services: map[types.NamespacedName]*apiv1.Service{
				svcNsName: createSvc("hi", 80),
			},
			resolver: alwaysTrueRefGrantResolver,
			name:     "BackendRef group wrong",
		},
		{
			gtr: wrongBackendRefKindGtr,
			expected: &L4Route{
				Source:     wrongBackendRefKindGtr,
				ParentRefs: []ParentRef{parentRefGraph},
				Spec: L4RouteSpec{
					Hostnames: []gatewayv1.Hostname{
						"app.example.com",
					},
					BackendRef: BackendRef{
						Valid:              false,
						InvalidForGateways: map[types.NamespacedName]conditions.Condition{},
					},
				},
				Conditions: []conditions.Condition{conditions.NewRouteBackendRefInvalidKind(
					"spec.rules[0].backendRefs[0].kind:" +
						" Unsupported value: \"not service\": supported values: \"Service\"",
				)},
				Attachable: true,
				Valid:      true,
			},
			gateway: createGateway(),
			services: map[types.NamespacedName]*apiv1.Service{
				svcNsName: createSvc("hi", 80),
			},
			resolver: alwaysTrueRefGrantResolver,
			name:     "BackendRef kind wrong",
		},
		{
			gtr: diffNsBackendRef,
			expected: &L4Route{
				Source:     diffNsBackendRef,
				ParentRefs: []ParentRef{parentRefGraph},
				Spec: L4RouteSpec{
					Hostnames: []gatewayv1.Hostname{
						"app.example.com",
					},
					BackendRef: BackendRef{
						Valid:              false,
						InvalidForGateways: map[types.NamespacedName]conditions.Condition{},
					},
				},
				Conditions: []conditions.Condition{conditions.NewRouteBackendRefRefNotPermitted(
					"spec.rules[0].backendRefs[0].namespace: Forbidden: Backend ref to Service " +
						"diff/hi not permitted by any ReferenceGrant",
				)},
				Attachable: true,
				Valid:      true,
			},
			gateway: createGateway(),
			services: map[types.NamespacedName]*apiv1.Service{
				diffSvcNsName: diffNsSvc,
			},
			resolver: alwaysFalseRefGrantResolver,
			name:     "BackendRef in diff namespace not permitted by any reference grant",
		},
		{
			gtr: portNilBackendRefGtr,
			expected: &L4Route{
				Source:     portNilBackendRefGtr,
				ParentRefs: []ParentRef{parentRefGraph},
				Spec: L4RouteSpec{
					Hostnames: []gatewayv1.Hostname{
						"app.example.com",
					},
					BackendRef: BackendRef{
						Valid:              false,
						InvalidForGateways: map[types.NamespacedName]conditions.Condition{},
					},
				},
				Conditions: []conditions.Condition{conditions.NewRouteBackendRefUnsupportedValue(
					"spec.rules[0].backendRefs[0].port: Required value: port cannot be nil",
				)},
				Attachable: true,
				Valid:      true,
			},
			gateway: createGateway(),
			services: map[types.NamespacedName]*apiv1.Service{
				diffSvcNsName: createSvc("hi", 80),
			},
			resolver: alwaysTrueRefGrantResolver,
			name:     "BackendRef port nil",
		},
		{
			gtr: ipFamilyMismatchGtr,
			expected: &L4Route{
				Source: ipFamilyMismatchGtr,
				ParentRefs: []ParentRef{
					{
						SectionName: helpers.GetPointer[gatewayv1.SectionName]("l1"),
						Gateway: &ParentRefGateway{
							NamespacedName: types.NamespacedName{
								Namespace: "test",
								Name:      "gateway",
							},
							EffectiveNginxProxy: &EffectiveNginxProxy{IPFamily: helpers.GetPointer(ngfAPI.IPv6)},
						},
					},
				},
				Spec: L4RouteSpec{
					Hostnames: []gatewayv1.Hostname{
						"app.example.com",
					},
					BackendRef: BackendRef{
						SvcNsName:   svcNsName,
						ServicePort: apiv1.ServicePort{Port: 80},
						InvalidForGateways: map[types.NamespacedName]conditions.Condition{
							{Namespace: "test", Name: "gateway"}: conditions.NewRouteInvalidIPFamily(
								"service configured with IPv4 family but NginxProxy is configured with IPv6",
							),
						},
						Valid: true,
					},
				},
				Attachable: true,
				Valid:      true,
			},
			gateway: modGateway(createGateway(), func(gw *Gateway) *Gateway {
				gw.EffectiveNginxProxy = &EffectiveNginxProxy{IPFamily: helpers.GetPointer(ngfAPI.IPv6)}
				return gw
			}),
			services: map[types.NamespacedName]*apiv1.Service{
				svcNsName: ipv4Svc,
			},
			resolver: alwaysTrueRefGrantResolver,
			name:     "service and npcfg ip family mismatch",
		},
		{
			gtr: diffNsBackendRef,
			expected: &L4Route{
				Source:     diffNsBackendRef,
				ParentRefs: []ParentRef{parentRefGraph},
				Spec: L4RouteSpec{
					Hostnames: []gatewayv1.Hostname{
						"app.example.com",
					},
					BackendRef: BackendRef{
						SvcNsName:          diffSvcNsName,
						ServicePort:        apiv1.ServicePort{Port: 80},
						Valid:              true,
						InvalidForGateways: map[types.NamespacedName]conditions.Condition{},
					},
				},
				Attachable: true,
				Valid:      true,
			},
			gateway: createGateway(),
			services: map[types.NamespacedName]*apiv1.Service{
				diffSvcNsName: diffNsSvc,
			},
			resolver: alwaysTrueRefGrantResolver,
			name:     "valid; backendRef in diff namespace permitted by a reference grant",
		},
		{
			gtr: validRefSameNs,
			expected: &L4Route{
				Source:     validRefSameNs,
				ParentRefs: []ParentRef{parentRefGraph},
				Spec: L4RouteSpec{
					Hostnames: []gatewayv1.Hostname{
						"app.example.com",
					},
					BackendRef: BackendRef{
						SvcNsName:          svcNsName,
						ServicePort:        apiv1.ServicePort{Port: 80},
						Valid:              true,
						InvalidForGateways: map[types.NamespacedName]conditions.Condition{},
					},
				},
				Attachable: true,
				Valid:      true,
			},
			gateway: createGateway(),
			services: map[types.NamespacedName]*apiv1.Service{
				svcNsName: ipv4Svc,
			},
			resolver: alwaysTrueRefGrantResolver,
			name:     "valid; same namespace",
		},
		{
			gtr: validRefSameNs,
			expected: &L4Route{
				Source:     validRefSameNs,
				ParentRefs: []ParentRef{parentRefGraph},
				Spec: L4RouteSpec{
					Hostnames: []gatewayv1.Hostname{
						"app.example.com",
					},
					BackendRef: BackendRef{
						SvcNsName:          svcNsName,
						ServicePort:        apiv1.ServicePort{Port: 80, AppProtocol: helpers.GetPointer(AppProtocolTypeWSS)},
						Valid:              true,
						InvalidForGateways: map[types.NamespacedName]conditions.Condition{},
					},
				},
				Attachable: true,
				Valid:      true,
			},
			gateway: createGateway(),
			services: map[types.NamespacedName]*apiv1.Service{
				svcNsName: createSvcWithAppProtocol("hi", AppProtocolTypeWSS, 80),
			},
			resolver: alwaysTrueRefGrantResolver,
			name:     "valid; same namespace, valid appProtocol",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			g := NewWithT(t)
			t.Parallel()

			r := buildTLSRoute(
				test.gtr,
				map[types.NamespacedName]*Gateway{client.ObjectKeyFromObject(test.gateway.Source): test.gateway},
				test.services,
				test.resolver,
			)
			g.Expect(helpers.Diff(test.expected, r)).To(BeEmpty())
		})
	}
}
