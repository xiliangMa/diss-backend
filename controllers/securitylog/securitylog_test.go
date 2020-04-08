package securitypolicy

import (
	"strings"
	"testing"
)

func Test_StringsTrim(t *testing.T) {
	t.Log(strings.TrimRight("abba", "ba"))
	t.Log(strings.TrimRight("abcdaaaaa", "abcd"))
	t.Log(strings.TrimSuffix(" select * from docker_ids where container_id = host and ", "and"))
}
