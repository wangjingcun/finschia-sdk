package keeper_test

import (
	"github.com/line/lbm-sdk/codec"
	"github.com/line/lbm-sdk/simapp"
	"github.com/line/lbm-sdk/testutil"
	sdk "github.com/line/lbm-sdk/types"
	paramskeeper "github.com/line/lbm-sdk/x/params/keeper"
)

func testComponents() (*codec.LegacyAmino, sdk.Context, sdk.StoreKey, paramskeeper.Keeper) {
	marshaler := simapp.MakeTestEncodingConfig().Marshaler
	legacyAmino := createTestCodec()
	mkey := sdk.NewKVStoreKey("test")
	tkey := sdk.NewTransientStoreKey("transient_test")
	ctx := testutil.DefaultContext(mkey, tkey)
	keeper := paramskeeper.NewKeeper(marshaler, legacyAmino, mkey, tkey)

	return legacyAmino, ctx, mkey, keeper
}

type invalid struct{}

type s struct {
	I int
}

func createTestCodec() *codec.LegacyAmino {
	cdc := codec.NewLegacyAmino()
	sdk.RegisterLegacyAminoCodec(cdc)
	cdc.RegisterConcrete(s{}, "test/s", nil)
	cdc.RegisterConcrete(invalid{}, "test/invalid", nil)
	return cdc
}
