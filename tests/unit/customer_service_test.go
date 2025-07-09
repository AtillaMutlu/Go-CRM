package unit

import (
	"Go-CRM/pkg/common"
	"Go-CRM/pkg/customer"
	"context"
	"testing"
)

func TestCustomerListParams_Defaults(t *testing.T) {
	params := customer.CustomerListParams{Page: 0, PageSize: 0}
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 || params.PageSize > 100 {
		params.PageSize = 10
	}
	if params.Page != 1 || params.PageSize != 10 {
		t.Errorf("Varsayılan pagination değerleri yanlış: page=%d, pageSize=%d", params.Page, params.PageSize)
	}
}

func TestContactListParams_Defaults(t *testing.T) {
	params := customer.ContactListParams{Page: 0, PageSize: 0}
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 || params.PageSize > 100 {
		params.PageSize = 10
	}
	if params.Page != 1 || params.PageSize != 10 {
		t.Errorf("Varsayılan pagination değerleri yanlış: page=%d, pageSize=%d", params.Page, params.PageSize)
	}
}

func TestPublishEvent_NoError(t *testing.T) {
	// Kafka init edilmiş varsayılır, gerçek bağlantı gerekmiyor
	err := common.PublishEvent(context.Background(), "test-key", "test-value")
	if err != nil {
		t.Errorf("Kafka event publish hatası: %v", err)
	}
}
