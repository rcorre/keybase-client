// Auto-generated by avdl-compiler v1.3.22 (https://github.com/keybase/node-avdl-compiler)
//   Input file: avdl/stellar1/local.avdl

package stellar1

import (
	"github.com/keybase/go-framed-msgpack-rpc/rpc"
	context "golang.org/x/net/context"
)

type DumpRes struct {
	Seed    string `codec:"seed" json:"seed"`
	Address string `codec:"address" json:"address"`
}

func (o DumpRes) DeepCopy() DumpRes {
	return DumpRes{
		Seed:    o.Seed,
		Address: o.Address,
	}
}

type BalancesLocalArg struct {
	AccountID AccountID `codec:"accountID" json:"accountID"`
}

type SendLocalArg struct {
	Recipient string `codec:"recipient" json:"recipient"`
	Amount    string `codec:"amount" json:"amount"`
	Asset     Asset  `codec:"asset" json:"asset"`
	Note      string `codec:"note" json:"note"`
}

type WalletInitArg struct {
}

type WalletDumpArg struct {
}

type LocalInterface interface {
	BalancesLocal(context.Context, AccountID) ([]Balance, error)
	SendLocal(context.Context, SendLocalArg) (PaymentResult, error)
	WalletInit(context.Context) error
	WalletDump(context.Context) (DumpRes, error)
}

func LocalProtocol(i LocalInterface) rpc.Protocol {
	return rpc.Protocol{
		Name: "stellar.1.local",
		Methods: map[string]rpc.ServeHandlerDescription{
			"balancesLocal": {
				MakeArg: func() interface{} {
					ret := make([]BalancesLocalArg, 1)
					return &ret
				},
				Handler: func(ctx context.Context, args interface{}) (ret interface{}, err error) {
					typedArgs, ok := args.(*[]BalancesLocalArg)
					if !ok {
						err = rpc.NewTypeError((*[]BalancesLocalArg)(nil), args)
						return
					}
					ret, err = i.BalancesLocal(ctx, (*typedArgs)[0].AccountID)
					return
				},
				MethodType: rpc.MethodCall,
			},
			"sendLocal": {
				MakeArg: func() interface{} {
					ret := make([]SendLocalArg, 1)
					return &ret
				},
				Handler: func(ctx context.Context, args interface{}) (ret interface{}, err error) {
					typedArgs, ok := args.(*[]SendLocalArg)
					if !ok {
						err = rpc.NewTypeError((*[]SendLocalArg)(nil), args)
						return
					}
					ret, err = i.SendLocal(ctx, (*typedArgs)[0])
					return
				},
				MethodType: rpc.MethodCall,
			},
			"walletInit": {
				MakeArg: func() interface{} {
					ret := make([]WalletInitArg, 1)
					return &ret
				},
				Handler: func(ctx context.Context, args interface{}) (ret interface{}, err error) {
					err = i.WalletInit(ctx)
					return
				},
				MethodType: rpc.MethodCall,
			},
			"walletDump": {
				MakeArg: func() interface{} {
					ret := make([]WalletDumpArg, 1)
					return &ret
				},
				Handler: func(ctx context.Context, args interface{}) (ret interface{}, err error) {
					ret, err = i.WalletDump(ctx)
					return
				},
				MethodType: rpc.MethodCall,
			},
		},
	}
}

type LocalClient struct {
	Cli rpc.GenericClient
}

func (c LocalClient) BalancesLocal(ctx context.Context, accountID AccountID) (res []Balance, err error) {
	__arg := BalancesLocalArg{AccountID: accountID}
	err = c.Cli.Call(ctx, "stellar.1.local.balancesLocal", []interface{}{__arg}, &res)
	return
}

func (c LocalClient) SendLocal(ctx context.Context, __arg SendLocalArg) (res PaymentResult, err error) {
	err = c.Cli.Call(ctx, "stellar.1.local.sendLocal", []interface{}{__arg}, &res)
	return
}

func (c LocalClient) WalletInit(ctx context.Context) (err error) {
	err = c.Cli.Call(ctx, "stellar.1.local.walletInit", []interface{}{WalletInitArg{}}, nil)
	return
}

func (c LocalClient) WalletDump(ctx context.Context) (res DumpRes, err error) {
	err = c.Cli.Call(ctx, "stellar.1.local.walletDump", []interface{}{WalletDumpArg{}}, &res)
	return
}
