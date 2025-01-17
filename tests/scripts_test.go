package tests

import (
	"fmt"
	"strings"
)

func subTestScripts(g *testGroup) {
	g.regSubTest("BASIC", scripts_BASIC_test)
	g.regSubTest("ATOMIC", scripts_ATOMIC_test)
	g.regSubTest("READONLY", scripts_READONLY_test)
	g.regSubTest("NONATOMIC", scripts_NONATOMIC_test)
}

func scripts_BASIC_test(mc *mockServer) error {
	return mc.DoBatch([][]interface{}{
		{"EVAL", "return 2 + 2", 0}, {"4"},
		{"SCRIPT LOAD", "return 2 + 2"}, {"2dd1b44209ecb49617af05caf0491390a03c1cc4"},
		{"SCRIPT EXISTS", "2dd1b44209ecb49617af05caf0491390a03c1cc4", "no_script"}, {"[1 0]"},
		{"EVALSHA", "2dd1b44209ecb49617af05caf0491390a03c1cc4", "0"}, {"4"},
		{"SCRIPT FLUSH"}, {"OK"},
		{"SCRIPT EXISTS", "2dd1b44209ecb49617af05caf0491390a03c1cc4", "no_script"}, {"[0 0]"},
		{"EVAL", "return KEYS[1] .. ' only'", 1, "key1"}, {"key1 only"},
		{"EVAL", "return KEYS[1] .. ' and ' .. ARGV[1]", 1, "key1", "arg1"}, {"key1 and arg1"},
		{"EVAL", "return ARGV[1] .. ' and ' .. ARGV[2]", 0, "arg1", "arg2"}, {"arg1 and arg2"},
		{"EVAL", "return tile38.sha1hex('asdf')", 0}, {"3da541559918a808c2402bba5012f6c60b27661c"},
		{"EVAL", "return tile38.distance_to(37.7341129, -122.4408378, 37.733, -122.43)", 0}, {"961"},
	})
}

func scripts_ATOMIC_test(mc *mockServer) error {
	return mc.DoBatch([][]interface{}{
		{"EVAL", "return tile38.call('get', KEYS[1], ARGV[1])", "1", "mykey", "myid"}, {nil},
		{"EVAL", "return tile38.call('set', KEYS[1], ARGV[1], 'point', 33, -115)", "1", "mykey", "myid1"}, {"OK"},
		{"EVAL", "return tile38.call('get', KEYS[1], ARGV[1], ARGV[2])", "1", "mykey", "myid1", "point"}, {"[33 -115]"},
	})
}

func scripts_READONLY_test(mc *mockServer) error {
	return mc.DoBatch([][]interface{}{
		{"EVALRO", "return tile38.call('get', KEYS[1], ARGV[1])", "1", "mykey", "myid"}, {nil},
		{"EVALRO", "return tile38.call('set', KEYS[1], ARGV[1], 'point', 33, -115)", "1", "mykey", "myid1"}, {
			func(v interface{}) (resp, expect interface{}) {
				s := fmt.Sprintf("%v", v)
				if strings.Contains(s, "ERR read only") {
					return v, v
				}
				return v, "A lua stack containing 'ERR read only'"
			},
		},
		{"EVALRO", "return tile38.pcall('set', KEYS[1], ARGV[1], 'point', 33, -115)", "1", "mykey", "myid1"}, {"ERR read only"},
		{"SET", "mykey", "myid1", "POINT", 33, -115}, {"OK"},
		{"EVALRO", "return tile38.call('get', KEYS[1], ARGV[1], ARGV[2])", "1", "mykey", "myid1", "point"}, {"[33 -115]"},
	})
}

func scripts_NONATOMIC_test(mc *mockServer) error {
	return mc.DoBatch([][]interface{}{
		{"EVALNA", "return tile38.call('get', KEYS[1], ARGV[1])", "1", "mykey", "myid"}, {nil},
		{"EVALNA", "return tile38.call('set', KEYS[1], ARGV[1], 'point', 33, -115)", "1", "mykey", "myid1"}, {"OK"},
		{"EVALNA", "return tile38.call('get', KEYS[1], ARGV[1], ARGV[2])", "1", "mykey", "myid1", "point"}, {"[33 -115]"},
	})
}
