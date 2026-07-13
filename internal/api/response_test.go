package api

import (
	"testing"
)

func TestErrorResponseJSONTags(t *testing.T) {
	errResp := ErrorResponse{
		Code:      "test_code",
		Message:   "test message",
		RequestID: "test_req_id",
	}
	if errResp.Code != "test_code" {
		t.Error("Code field mismatch")
	}
	if errResp.Message != "test message" {
		t.Error("Message field mismatch")
	}
	if errResp.RequestID != "test_req_id" {
		t.Error("RequestID field mismatch")
	}
}
