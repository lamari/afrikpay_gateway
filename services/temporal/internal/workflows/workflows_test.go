package workflows

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.temporal.io/sdk/testsuite"
	"github.com/afrikpay/gateway/internal/model"
	"github.com/afrikpay/gateway/internal/activities"
	"github.com/afrikpay/gateway/internal/config"
)

func TestCreateUserWorkflow(t *testing.T) {
	suite := &testsuite.WorkflowTestSuite{}
	env := suite.NewTestWorkflowEnvironment()

	// Mock activity result
	env.RegisterActivity(activities.CreateUserActivity)
	env.OnActivity("CreateUserActivity", mock.Anything, mock.Anything).
		Return("mocked-uid", nil)

	var result string
	input := model.UserInput{Firstname: "Alaa", Lastname: "Lamari", Email: "a@b.com"}
	cfg := &config.Config{}
	SetConfig(cfg)
	env.ExecuteWorkflow(CreateUserWorkflow, input)
	assert.True(t, env.IsWorkflowCompleted())
	assert.NoError(t, env.GetWorkflowError())
	env.GetWorkflowResult(&result)
	assert.Equal(t, "mocked-uid", result)
}

func TestCreateWalletWorkflow(t *testing.T) {
	suite := &testsuite.WorkflowTestSuite{}
	env := suite.NewTestWorkflowEnvironment()
	env.RegisterActivity(activities.CreateWalletActivity)
	env.OnActivity("CreateWalletActivity", mock.Anything, mock.Anything).
		Return("mocked-wallet", nil)

	var result string
	input := model.WalletInput{UserUID: "u1", Currency: "XAF"}
	cfg := &config.Config{}
	SetConfig(cfg)
	env.ExecuteWorkflow(CreateWalletWorkflow, input)
	assert.True(t, env.IsWorkflowCompleted())
	assert.NoError(t, env.GetWorkflowError())
	env.GetWorkflowResult(&result)
	assert.Equal(t, "mocked-wallet", result)
}

func TestDepositMobileMoneyWorkflow(t *testing.T) {
	suite := &testsuite.WorkflowTestSuite{}
	env := suite.NewTestWorkflowEnvironment()
	env.RegisterActivity(activities.DepositMobileMoneyActivity)
	env.OnActivity("DepositMobileMoneyActivity", mock.Anything, mock.Anything).
		Return("success", nil)

	var result string
	input := model.DepositInput{UserUID: "u1", Amount: 1000, Provider: "MTN", Phone: "+237600000000"}
	cfg := &config.Config{}
	SetConfig(cfg)
	env.ExecuteWorkflow(DepositMobileMoneyWorkflow, input)
	assert.True(t, env.IsWorkflowCompleted())
	assert.NoError(t, env.GetWorkflowError())
	env.GetWorkflowResult(&result)
	assert.Equal(t, "success", result)
}

func TestBuyCryptoWorkflow(t *testing.T) {
	suite := &testsuite.WorkflowTestSuite{}
	env := suite.NewTestWorkflowEnvironment()
	env.RegisterActivity(activities.BuyCryptoActivity)
	env.OnActivity("BuyCryptoActivity", mock.Anything, mock.Anything).
		Return("success", nil)

	var result string
	input := model.BuyCryptoInput{UserUID: "u1", Symbol: "BTC", Amount: 100}
	cfg := &config.Config{}
	SetConfig(cfg)
	env.ExecuteWorkflow(BuyCryptoWorkflow, input)
	assert.True(t, env.IsWorkflowCompleted())
	assert.NoError(t, env.GetWorkflowError())
	env.GetWorkflowResult(&result)
	assert.Equal(t, "success", result)
}

func TestSubmitKYCWorkflow(t *testing.T) {
	suite := &testsuite.WorkflowTestSuite{}
	env := suite.NewTestWorkflowEnvironment()
	env.RegisterActivity(activities.SubmitKYCActivity)
	env.OnActivity("SubmitKYCActivity", mock.Anything, mock.Anything).
		Return("kyc-ok", nil)

	var result string
	input := model.KYCInput{UserUID: "u1", DocumentType: "CNI", DocumentNumber: "123456"}
	cfg := &config.Config{}
	SetConfig(cfg)
	env.ExecuteWorkflow(SubmitKYCWorkflow, input)
	assert.True(t, env.IsWorkflowCompleted())
	assert.NoError(t, env.GetWorkflowError())
	env.GetWorkflowResult(&result)
	assert.Equal(t, "kyc-ok", result)
}
