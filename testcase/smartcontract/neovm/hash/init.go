package hash

import (
	"github.com/ontio/ontology-test/testframework"
)

func TestHash() {
	testframework.TFramework.RegTestCase("TestHash160", TestHash160)
	testframework.TFramework.RegTestCase("TestHash256", TestHash256)
	testframework.TFramework.RegTestCase("TestSha1", TestSha1)
	testframework.TFramework.RegTestCase("TestSha256", TestSha256)
}
