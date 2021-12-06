package main

import "github.com/kanyuanzhi/tialloy-client/ticnet"

func main(){
	client := ticnet.NewClient()
	client.Serve()
}
