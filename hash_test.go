package ovoclient

import(
	"testing"
)

func TestConfigurationLoad(t *testing.T) {
	t.Logf("Test key = %s -> Hash = %d", "test12345", GetPositiveHashCode("test12345",128))
	t.Logf("Test key = %s -> Hash = %d", "ciaociao", GetPositiveHashCode("ciaociao",128))
	t.Logf("Test key = %s -> Hash = %d", "asdfghjklòàèé", GetPositiveHashCode("asdfghjklòàèé",128))
	t.Logf("Test key = %s -> Hash = %d", "你好 你好 你好", GetPositiveHashCode("你好 你好 你好",128))
	t.Logf("Test key = %s -> Hash = %d", "cammello", GetPositiveHashCode("cammello",128))
	t.Logf("Test key = %s -> Hash = %d", "早上好，女士们", GetPositiveHashCode("早上好，女士们",128))
	t.Logf("Test key = %s -> Hash = %d", "èéà", GetPositiveHashCode("èéà",128))
	
}