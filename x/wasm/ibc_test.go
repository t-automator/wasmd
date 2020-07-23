package wasm_test

import (
	"testing"

	"github.com/CosmWasm/wasmd/x/wasm"
	sdk "github.com/cosmos/cosmos-sdk/types"
	connectiontypes "github.com/cosmos/cosmos-sdk/x/ibc/03-connection/types"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/types"
	host "github.com/cosmos/cosmos-sdk/x/ibc/24-host"
	"github.com/stretchr/testify/suite"
)

// define constants used for testing
const (
	testClientIDA = "testclientIDA"
	testClientIDB = "testClientIDb"

	testConnection = "testconnectionatob"
	testPort1      = "ibc-wasm"
	testPort2      = "testportid"
	testChannel1   = "firstchannel"
	testChannel2   = "secondchannel"
)

// define variables used for testing
var (
	testAddr1, _ = sdk.AccAddressFromBech32("cosmos1scqhwpgsmr6vmztaa7suurfl52my6nd2kmrudl")
	testAddr2, _ = sdk.AccAddressFromBech32("cosmos1scqhwpgsmr6vmztaa7suurfl52my6nd2kmrujl")

	testCoins, _ = sdk.ParseCoins("100atom")
	prefixCoins  = sdk.NewCoins(sdk.NewCoin("bank/firstchannel/atom", sdk.NewInt(100)))
	prefixCoins2 = sdk.NewCoins(sdk.NewCoin("testportid/secondchannel/atom", sdk.NewInt(100)))
)

type IBCTestSuite struct {
	suite.Suite

	chainA *TestChain
	chainB *TestChain

	cleanupA func()
	cleanupB func()
}

func (suite *IBCTestSuite) SetupTest() {
	suite.chainA, suite.cleanupA = NewTestChain(testClientIDA)
	suite.chainB, suite.cleanupB = NewTestChain(testClientIDB)
}

func (suite *IBCTestSuite) TearDownTest() {
	suite.cleanupA()
	suite.cleanupB()
}

func (suite *IBCTestSuite) TestBindPorts() {
	suite.T().Logf("To be implemented")
}

func (suite *IBCTestSuite) TestHandleMsgTransfer() {
	// create channel capability from ibc scoped keeper and claim with transfer scoped keeper
	capName := host.ChannelCapabilityPath(testPort1, testChannel1)
	cap, err := suite.chainA.App.ScopedIBCKeeper.NewCapability(suite.chainA.GetContext(), capName)
	suite.Require().NoError(err)
	err = suite.chainA.App.ScopedWasmKeeper.ClaimCapability(suite.chainA.GetContext(), cap, capName)
	suite.Require().NoError(err)

	handler := wasm.NewHandler(suite.chainA.App.WasmKeeper)

	var (
		destContractAddr, _ = sdk.AccAddressFromBech32("cosmos18vd8fpwxzck93qlwghaj6arh4p7c5n89uzcee5")
	)

	ctx := suite.chainA.GetContext()
	msg := &wasm.MsgWasmIBCCall{
		SourcePort:       testPort1,
		SourceChannel:    testChannel1,
		Sender:           testAddr1,
		DestContractAddr: destContractAddr,
		TimeoutHeight:    110,
		TimeoutTimestamp: 0,
		Msg:              []byte("{}"),
	}
	// Setup channel from A to B
	suite.chainA.CreateClient(suite.chainB)
	suite.chainA.createConnection(testConnection, testConnection, testClientIDB, testClientIDA, connectiontypes.OPEN)
	suite.chainA.createChannel(testPort1, testChannel1, testPort2, testChannel2, channeltypes.OPEN, channeltypes.ORDERED, testConnection)

	nextSeqSend := uint64(1)
	suite.chainA.App.IBCKeeper.ChannelKeeper.SetNextSequenceSend(ctx, testPort1, testChannel1, nextSeqSend)

	_ = suite.chainA.App.BankKeeper.SetBalances(ctx, testAddr1, testCoins)
	res, err := handler(ctx, msg)
	suite.Require().NoError(err)
	suite.Require().NotNil(res, "%+v", res) // successfully executed

}

func TestIBCTestSuite(t *testing.T) {
	suite.Run(t, new(IBCTestSuite))
}
