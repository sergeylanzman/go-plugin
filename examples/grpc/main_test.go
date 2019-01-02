package main

import (
	"fmt"
	"os"
	"os/exec"
	"plugin"
	"testing"

	"github.com/hashicorp/go-plugin/examples/grpc/native"

	"github.com/hashicorp/go-hclog"
	hplugin "github.com/hashicorp/go-plugin"
	"github.com/hashicorp/go-plugin/examples/grpc/shared"
)

var grpcKV, clientGrpc = getPlugin("kv_grpc", "./kv-go-grpc")
var grpcPython, clientPyhton = getPlugin("kv_grpc", "python plugin-python/plugin.py")
var netRpcKV, clientNetRpc = getPlugin("kv", "./kv-go-netrpc")

func getPlugin(pluginName, kvPlugin string) (shared.KV, *hplugin.Client) {
	client := hplugin.NewClient(&hplugin.ClientConfig{
		HandshakeConfig: shared.Handshake,
		Plugins:         shared.PluginMap,
		Managed:         true,
		Cmd:             exec.Command("sh", "-c", kvPlugin),
		AllowedProtocols: []hplugin.Protocol{
			hplugin.ProtocolNetRPC, hplugin.ProtocolGRPC},
		Logger: hclog.New(&hclog.LoggerOptions{
			Name:   "plugin",
			Output: os.Stdout,
			Level:  hclog.Info,
		}),
	})

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	// Request the plugin
	raw, err := rpcClient.Dispense(pluginName)
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	kv := raw.(shared.KV)
	return kv, client
}

func BenchmarkGrpc(b *testing.B) {
	for n := 0; n < b.N; n++ {
		grpcKV.Bench()
	}
	b.StopTimer()
	clientGrpc.Kill()
}

func BenchmarkGrpcPyhton(b *testing.B) {
	for n := 0; n < b.N; n++ {
		grpcPython.Bench()
	}
	b.StopTimer()
	clientPyhton.Kill()
}

func BenchmarkNetRpc(b *testing.B) {
	for n := 0; n < b.N; n++ {
		netRpcKV.Bench()
	}
	b.StopTimer()
	clientNetRpc.Kill()
}

func BenchmarkNative(b *testing.B) {
	b.StopTimer()
	kv := native.KVType{}
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		kv.Bench()
	}
}

func BenchmarkPlugin(b *testing.B) {
	b.StopTimer()
	plug, err := plugin.Open("native.so")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	p, err := plug.Lookup("KV")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	kv, ok := p.(shared.KV)
	if !ok {
		fmt.Println("unexpected type from module symbol")
		os.Exit(1)
	}
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		kv.Bench()
	}
}
