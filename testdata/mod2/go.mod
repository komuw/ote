module testdata/mod2

go 1.16

replace github.com/shirou/gopsutil => github.com/hashicorp/gopsutil v2.18.13-0.20200531184148-5aca383d4f9d+incompatible

require (
	github.com/LK4D4/joincontext v0.0.0-20171026170139-1724345da6d5
	github.com/StackExchange/wmi v0.0.0-20210224194228-fe8f1750fd46 // indirect
	github.com/frankban/quicktest v1.12.1
	github.com/go-ole/go-ole v1.2.5 // indirect
	github.com/hashicorp/nomad v1.0.4
	github.com/pkg/errors v0.9.1
	github.com/shirou/gopsutil v2.20.9+incompatible
	github.com/stretchr/testify v1.7.0 // indirect
	golang.org/x/net v0.0.0-20210501142056-aec3718b3fa0 // indirect
	golang.org/x/sys v0.0.0-20210426230700-d19ff857e887 // indirect
	rsc.io/quote v1.5.2
)
