package auth

import (
	"errors"
	"fmt"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	_ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	kexec "k8s.io/client-go/plugin/pkg/client/auth/exec"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	_ "k8s.io/client-go/plugin/pkg/client/auth/openstack"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/transport/spdy"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"net/http"
)

func GetFactory() (cmdutil.Factory){
	kubeConfigFlags := genericclioptions.NewConfigFlags(true).WithDeprecatedPasswordFlag()
	matchVersionKubeConfigFlags := cmdutil.NewMatchVersionFlags(kubeConfigFlags)
	return cmdutil.NewFactory(matchVersionKubeConfigFlags)

}

func GetExecAuthToken() (string, error){
	f:= GetFactory()
	r,_ := f.ToRESTConfig()
	if r.ExecProvider == nil {
		return "", errors.New("ExecProvider is not found")
	}
	a,_:= kexec.GetAuthenticator(r.ExecProvider)
	rt,_,_ := spdy.RoundTripperFor(r)
	t, _ :=r.TransportConfig()

	a.UpdateTransportConfig(t)
	ht,_ := http.NewRequest("GET",r.Host,nil)
	t.WrapTransport(rt).RoundTrip(ht)
	token := ht.Header.Get("Authorization")
	if token !="" {
		return token, nil
	}
	return "", nil
}

func RefreshAuthToken() (bool){
	f:= GetFactory()
	r,_ := f.ToRESTConfig()

	provider, err := rest.GetAuthProvider(r.Host, r.AuthProvider, r.AuthConfigPersister)

	if err != nil {
		fmt.Print(err)
	}
	err = provider.Login()
	if err != nil {
		fmt.Print(err)
	}
	rt,_,_ := spdy.RoundTripperFor(r)

	ht,_ := http.NewRequest("GET",r.Host,nil)
	rt1 := provider.WrapTransport(rt)
	rt1.RoundTrip(ht)
	return true

}


func main() {
	fmt.Println("Test")
	kubeConfigFlags := genericclioptions.NewConfigFlags(true).WithDeprecatedPasswordFlag()
	matchVersionKubeConfigFlags := cmdutil.NewMatchVersionFlags(kubeConfigFlags)
	f := cmdutil.NewFactory(matchVersionKubeConfigFlags)
	r,_ := f.ToRESTConfig()

	if r.AuthProvider != nil {

		provider, err := rest.GetAuthProvider(r.Host, r.AuthProvider, r.AuthConfigPersister)
		if err != nil {
			fmt.Print(err)
		}
		err = provider.Login()
		if err != nil {
			fmt.Print(err)
		}
		rt,_,_ := spdy.RoundTripperFor(r)

		ht,_ := http.NewRequest("GET",r.Host,nil)
		rt1 := provider.WrapTransport(rt)
		rt1.RoundTrip(ht)

		fmt.Println("Auth token",ht.Header.Get("Authorization"))
		return
	}
	//if r.ExecProvider == nil {
	//	return
	//}
	a,_:= kexec.GetAuthenticator(r.ExecProvider)
	rt,_,_ := spdy.RoundTripperFor(r)
	t, _ :=r.TransportConfig()

	a.UpdateTransportConfig(t)
	ht,_ := http.NewRequest("GET","",nil)
	t.WrapTransport(rt).RoundTrip(ht)
	//kexec.
	fmt.Println("token",ht.Header.Get("Authorization"))

}
