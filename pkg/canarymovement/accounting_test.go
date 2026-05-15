package canarymovement

import "testing"

func TestRealizedRequiresPostVerify(t *testing.T) {
	doc := validSurface()
	doc["postVerifyPassed"] = false
	if err := evaluateAccounting(doc); err == nil {
		t.Fatal("expected rejection when realized compression without post-verify")
	}
}
