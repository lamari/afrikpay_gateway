package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"context"
	goTemporalClient "go.temporal.io/sdk/client"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// --- Mocks ---
type mockWorkflowRun struct{
	result interface{}
	err    error
}
// Ajout pour compatibilité interface WorkflowRun
func (m *mockWorkflowRun) GetWithOptions(ctx context.Context, valuePtr interface{}, opts goTemporalClient.WorkflowRunGetOptions) error {
	return m.Get(ctx, valuePtr)
}
func (m *mockWorkflowRun) Get(ctx context.Context, valuePtr interface{}) error {
	if m.err != nil {
		return m.err
	}
	b, _ := json.Marshal(m.result)
	return json.Unmarshal(b, valuePtr)
}
// Les autres méthodes de l'interface WorkflowRun ne sont pas utilisées ici, on peut les laisser vides.
func (m *mockWorkflowRun) GetID() string { return "mock-id" }
func (m *mockWorkflowRun) GetRunID() string { return "mock-run-id" }
func (m *mockWorkflowRun) GetWorkflowType() string { return "mock-type" }
func (m *mockWorkflowRun) GetStartTime() int64 { return 0 }
func (m *mockWorkflowRun) GetCloseTime() int64 { return 0 }
func (m *mockWorkflowRun) GetStatus() int { return 0 }
func (m *mockWorkflowRun) GetHistoryLength() int64 { return 0 }
func (m *mockWorkflowRun) GetMemo() map[string]interface{} { return nil }
func (m *mockWorkflowRun) GetSearchAttributes() map[string]interface{} { return nil }
func (m *mockWorkflowRun) GetRawHistory() []byte { return nil }
func (m *mockWorkflowRun) GetTraceID() string { return "" }


type mockTemporalClient struct{
	executeFunc func(ctx context.Context, options goTemporalClient.StartWorkflowOptions, workflow interface{}, args ...interface{}) (goTemporalClient.WorkflowRun, error)
}
func (m *mockTemporalClient) ExecuteWorkflow(ctx context.Context, options goTemporalClient.StartWorkflowOptions, workflow interface{}, args ...interface{}) (goTemporalClient.WorkflowRun, error) {
	return m.executeFunc(ctx, options, workflow, args...)
}

// --- Tests ---
func TestWorkflowHandler_CreateUser(t *testing.T) {
	e := echo.New()
	input := map[string]interface{}{"email": "test@afrikpay.com", "password": "pass"}
	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/api/workflow/v1/CreateUser", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("version", "nameworkflow")
	c.SetParamValues("v1", "CreateUser")

	SetTemporalClient(&mockTemporalClient{
		executeFunc: func(ctx context.Context, options goTemporalClient.StartWorkflowOptions, workflow interface{}, args ...interface{}) (goTemporalClient.WorkflowRun, error) {
			return &mockWorkflowRun{result: "user-uid-123", err: nil}, nil
		},
	})

	if assert.NoError(t, WorkflowHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var resp map[string]string
		json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.Equal(t, "user-uid-123", resp["user_uid"])
	}
}

func TestWorkflowHandler_InvalidInput(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/workflow/v1/CreateUser", bytes.NewReader([]byte("notjson")))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("version", "nameworkflow")
	c.SetParamValues("v1", "CreateUser")

	SetTemporalClient(&mockTemporalClient{
		executeFunc: func(ctx context.Context, options goTemporalClient.StartWorkflowOptions, workflow interface{}, args ...interface{}) (goTemporalClient.WorkflowRun, error) {
			return &mockWorkflowRun{result: "", err: nil}, nil
		},
	})

	err := WorkflowHandler(c)
	assert.Error(t, err)
}

func TestWorkflowHandler_UnknownWorkflow(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/workflow/v1/Unknown", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("version", "nameworkflow")
	c.SetParamValues("v1", "Unknown")

	SetTemporalClient(&mockTemporalClient{
		executeFunc: func(ctx context.Context, options goTemporalClient.StartWorkflowOptions, workflow interface{}, args ...interface{}) (goTemporalClient.WorkflowRun, error) {
			return &mockWorkflowRun{result: "", err: nil}, nil
		},
	})

	_ = WorkflowHandler(c)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestWorkflowHandler_WrongVersion(t *testing.T) {
	e := echo.New()
	input := map[string]interface{}{"email": "test@afrikpay.com", "password": "pass"}
	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/api/workflow/v2/CreateUser", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("version", "nameworkflow")
	c.SetParamValues("v2", "CreateUser")

	SetTemporalClient(&mockTemporalClient{
		executeFunc: func(ctx context.Context, options goTemporalClient.StartWorkflowOptions, workflow interface{}, args ...interface{}) (goTemporalClient.WorkflowRun, error) {
			return &mockWorkflowRun{result: "user-uid-123", err: nil}, nil
		},
	})

	_ = WorkflowHandler(c)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestWorkflowHandler_WrongName(t *testing.T) {
	e := echo.New()
	input := map[string]interface{}{"email": "test@afrikpay.com", "password": "pass"}
	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/api/workflow/v1/UnknownWorkflow", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("version", "nameworkflow")
	c.SetParamValues("v1", "UnknownWorkflow")

	SetTemporalClient(&mockTemporalClient{
		executeFunc: func(ctx context.Context, options goTemporalClient.StartWorkflowOptions, workflow interface{}, args ...interface{}) (goTemporalClient.WorkflowRun, error) {
			return &mockWorkflowRun{result: "user-uid-123", err: nil}, nil
		},
	})

	_ = WorkflowHandler(c)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestWorkflowHandler_WrongNameAndVersion(t *testing.T) {
	e := echo.New()
	input := map[string]interface{}{"email": "test@afrikpay.com", "password": "pass"}
	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/api/workflow/v2/UnknownWorkflow", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("version", "nameworkflow")
	c.SetParamValues("v2", "UnknownWorkflow")

	SetTemporalClient(&mockTemporalClient{
		executeFunc: func(ctx context.Context, options goTemporalClient.StartWorkflowOptions, workflow interface{}, args ...interface{}) (goTemporalClient.WorkflowRun, error) {
			return &mockWorkflowRun{result: "user-uid-123", err: nil}, nil
		},
	})

	_ = WorkflowHandler(c)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestWorkflowHandler_NameVersion_Success(t *testing.T) {
	e := echo.New()
	input := map[string]interface{}{"email": "test@afrikpay.com", "password": "pass"}
	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/api/workflow/v1/CreateUser", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("version", "nameworkflow")
	c.SetParamValues("v1", "CreateUser")

	SetTemporalClient(&mockTemporalClient{
		executeFunc: func(ctx context.Context, options goTemporalClient.StartWorkflowOptions, workflow interface{}, args ...interface{}) (goTemporalClient.WorkflowRun, error) {
			return &mockWorkflowRun{result: "user-uid-123", err: nil}, nil
		},
	})

	if assert.NoError(t, WorkflowHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var resp map[string]string
		json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.Equal(t, "user-uid-123", resp["user_uid"])
	}
}


func TestWorkflowHandler_CreateWallet(t *testing.T) {
	e := echo.New()
	input := map[string]interface{}{"user_uid": "u1", "currency": "XAF"}
	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/api/workflow/v1/CreateWallet", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("version", "nameworkflow")
	c.SetParamValues("v1", "CreateWallet")

	SetTemporalClient(&mockTemporalClient{
		executeFunc: func(ctx context.Context, options goTemporalClient.StartWorkflowOptions, workflow interface{}, args ...interface{}) (goTemporalClient.WorkflowRun, error) {
			return &mockWorkflowRun{result: "wallet-uid-456", err: nil}, nil
		},
	})

	if assert.NoError(t, WorkflowHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var resp map[string]string
		json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.Equal(t, "wallet-uid-456", resp["wallet_uid"])
	}
}

func TestWorkflowHandler_CreateWallet_InvalidInput(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/workflow/v1/CreateWallet", bytes.NewReader([]byte("notjson")))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("version", "nameworkflow")
	c.SetParamValues("v1", "CreateWallet")

	SetTemporalClient(&mockTemporalClient{
		executeFunc: func(ctx context.Context, options goTemporalClient.StartWorkflowOptions, workflow interface{}, args ...interface{}) (goTemporalClient.WorkflowRun, error) {
			return &mockWorkflowRun{result: "", err: nil}, nil
		},
	})

	err := WorkflowHandler(c)
	assert.Error(t, err)
}

func TestWorkflowHandler_CreateWallet_TemporalError(t *testing.T) {
	e := echo.New()
	input := map[string]interface{}{ "user_uid": "u1", "currency": "XAF" }
	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/api/workflow/v1/CreateWallet", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("version", "nameworkflow")
	c.SetParamValues("v1", "CreateWallet")

	SetTemporalClient(&mockTemporalClient{
		executeFunc: func(ctx context.Context, options goTemporalClient.StartWorkflowOptions, workflow interface{}, args ...interface{}) (goTemporalClient.WorkflowRun, error) {
			return &mockWorkflowRun{result: "", err: nil}, assert.AnError
		},
	})

	err := WorkflowHandler(c)
	assert.Error(t, err)
}

// --- DepositMobileMoney ---
func TestWorkflowHandler_DepositMobileMoney(t *testing.T) {
	e := echo.New()
	input := map[string]interface{}{ "user_uid": "u1", "amount": 1000, "provider": "MTN", "phone": "123456789" }
	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/api/workflow/v1/DepositMobileMoney", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("version", "nameworkflow")
	c.SetParamValues("v1", "DepositMobileMoney")

	SetTemporalClient(&mockTemporalClient{
		executeFunc: func(ctx context.Context, options goTemporalClient.StartWorkflowOptions, workflow interface{}, args ...interface{}) (goTemporalClient.WorkflowRun, error) {
			return &mockWorkflowRun{result: "deposit-success", err: nil}, nil
		},
	})

	if assert.NoError(t, WorkflowHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var resp map[string]string
		json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.Equal(t, "deposit-success", resp["result"])
	}
}

func TestWorkflowHandler_DepositMobileMoney_InvalidInput(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/workflow/v1/DepositMobileMoney", bytes.NewReader([]byte("notjson")))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("version", "nameworkflow")
	c.SetParamValues("v1", "DepositMobileMoney")

	SetTemporalClient(&mockTemporalClient{
		executeFunc: func(ctx context.Context, options goTemporalClient.StartWorkflowOptions, workflow interface{}, args ...interface{}) (goTemporalClient.WorkflowRun, error) {
			return &mockWorkflowRun{result: "", err: nil}, nil
		},
	})

	err := WorkflowHandler(c)
	assert.Error(t, err)
}

func TestWorkflowHandler_DepositMobileMoney_TemporalError(t *testing.T) {
	e := echo.New()
	input := map[string]interface{}{ "user_uid": "u1", "amount": 1000, "provider": "MTN", "phone": "123456789" }
	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/api/workflow/v1/DepositMobileMoney", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("version", "nameworkflow")
	c.SetParamValues("v1", "DepositMobileMoney")

	SetTemporalClient(&mockTemporalClient{
		executeFunc: func(ctx context.Context, options goTemporalClient.StartWorkflowOptions, workflow interface{}, args ...interface{}) (goTemporalClient.WorkflowRun, error) {
			return &mockWorkflowRun{result: "", err: nil}, assert.AnError
		},
	})

	err := WorkflowHandler(c)
	assert.Error(t, err)
}

// --- BuyCrypto ---
func TestWorkflowHandler_BuyCrypto(t *testing.T) {
	e := echo.New()
	input := map[string]interface{}{ "user_uid": "u1", "symbol": "BTC", "amount": 0.01 }
	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/api/workflow/v1/BuyCrypto", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("version", "nameworkflow")
	c.SetParamValues("v1", "BuyCrypto")

	SetTemporalClient(&mockTemporalClient{
		executeFunc: func(ctx context.Context, options goTemporalClient.StartWorkflowOptions, workflow interface{}, args ...interface{}) (goTemporalClient.WorkflowRun, error) {
			return &mockWorkflowRun{result: "buy-success", err: nil}, nil
		},
	})

	if assert.NoError(t, WorkflowHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var resp map[string]string
		json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.Equal(t, "buy-success", resp["result"])
	}
}

func TestWorkflowHandler_BuyCrypto_InvalidInput(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/workflow/v1/BuyCrypto", bytes.NewReader([]byte("notjson")))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("version", "nameworkflow")
	c.SetParamValues("v1", "BuyCrypto")

	SetTemporalClient(&mockTemporalClient{
		executeFunc: func(ctx context.Context, options goTemporalClient.StartWorkflowOptions, workflow interface{}, args ...interface{}) (goTemporalClient.WorkflowRun, error) {
			return &mockWorkflowRun{result: "", err: nil}, nil
		},
	})

	err := WorkflowHandler(c)
	assert.Error(t, err)
}

func TestWorkflowHandler_BuyCrypto_TemporalError(t *testing.T) {
	e := echo.New()
	input := map[string]interface{}{ "user_uid": "u1", "symbol": "BTC", "amount": 0.01 }
	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/api/workflow/v1/BuyCrypto", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("version", "nameworkflow")
	c.SetParamValues("v1", "BuyCrypto")

	SetTemporalClient(&mockTemporalClient{
		executeFunc: func(ctx context.Context, options goTemporalClient.StartWorkflowOptions, workflow interface{}, args ...interface{}) (goTemporalClient.WorkflowRun, error) {
			return &mockWorkflowRun{result: "", err: nil}, assert.AnError
		},
	})

	err := WorkflowHandler(c)
	assert.Error(t, err)
}

// --- SubmitKYC ---
func TestWorkflowHandler_SubmitKYC(t *testing.T) {
	e := echo.New()
	input := map[string]interface{}{ "user_uid": "u1", "document_type": "ID", "document_number": "1234", "document_url": "url" }
	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/api/workflow/v1/SubmitKYC", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("version", "nameworkflow")
	c.SetParamValues("v1", "SubmitKYC")

	SetTemporalClient(&mockTemporalClient{
		executeFunc: func(ctx context.Context, options goTemporalClient.StartWorkflowOptions, workflow interface{}, args ...interface{}) (goTemporalClient.WorkflowRun, error) {
			return &mockWorkflowRun{result: "kyc-success", err: nil}, nil
		},
	})

	if assert.NoError(t, WorkflowHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var resp map[string]string
		json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.Equal(t, "kyc-success", resp["result"])
	}
}

func TestWorkflowHandler_SubmitKYC_InvalidInput(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/workflow/v1/SubmitKYC", bytes.NewReader([]byte("notjson")))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("version", "nameworkflow")
	c.SetParamValues("v1", "SubmitKYC")

	SetTemporalClient(&mockTemporalClient{
		executeFunc: func(ctx context.Context, options goTemporalClient.StartWorkflowOptions, workflow interface{}, args ...interface{}) (goTemporalClient.WorkflowRun, error) {
			return &mockWorkflowRun{result: "", err: nil}, nil
		},
	})

	err := WorkflowHandler(c)
	assert.Error(t, err)
}

func TestWorkflowHandler_SubmitKYC_TemporalError(t *testing.T) {
	e := echo.New()
	input := map[string]interface{}{ "user_uid": "u1", "document_type": "ID", "document_number": "1234", "document_url": "url" }
	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/api/workflow/v1/SubmitKYC", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("version", "nameworkflow")
	c.SetParamValues("v1", "SubmitKYC")

	SetTemporalClient(&mockTemporalClient{
		executeFunc: func(ctx context.Context, options goTemporalClient.StartWorkflowOptions, workflow interface{}, args ...interface{}) (goTemporalClient.WorkflowRun, error) {
			return &mockWorkflowRun{result: "", err: nil}, assert.AnError
		},
	})

	err := WorkflowHandler(c)
	assert.Error(t, err)
}
