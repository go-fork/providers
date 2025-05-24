package config

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManager_NewConfig(t *testing.T) {
	cfg := NewConfig()
	assert.NotNil(t, cfg, "NewConfig should return a non-nil Manager")
}

func TestManager_GetSet(t *testing.T) {
	cfg := NewConfig()

	// Test Set with empty key
	err := cfg.Set("", "value")
	assert.Error(t, err, "Set with empty key should return an error")

	// Test Set and Get with various data types
	tests := []struct {
		key      string
		value    interface{}
		getFunc  func(string) (interface{}, bool)
		expected interface{}
	}{
		{
			key:   "string.key",
			value: "string value",
			getFunc: func(key string) (interface{}, bool) {
				return cfg.GetString(key)
			},
			expected: "string value",
		},
		{
			key:   "int.key",
			value: 42,
			getFunc: func(key string) (interface{}, bool) {
				return cfg.GetInt(key)
			},
			expected: 42,
		},
		{
			key:   "bool.key",
			value: true,
			getFunc: func(key string) (interface{}, bool) {
				return cfg.GetBool(key)
			},
			expected: true,
		},
		{
			key:   "float.key",
			value: 3.14,
			getFunc: func(key string) (interface{}, bool) {
				return cfg.GetFloat(key)
			},
			expected: 3.14,
		},
		{
			key:   "duration.key",
			value: "5s",
			getFunc: func(key string) (interface{}, bool) {
				return cfg.GetDuration(key)
			},
			expected: 5 * time.Second,
		},
		{
			key:   "generic.key",
			value: map[string]interface{}{"nested": "value"},
			getFunc: func(key string) (interface{}, bool) {
				return cfg.Get(key)
			},
			expected: map[string]interface{}{"nested": "value"},
		},
	}

	for _, test := range tests {
		t.Run(test.key, func(t *testing.T) {
			err := cfg.Set(test.key, test.value)
			assert.NoError(t, err, "Set should not return an error")

			// Verify Has method
			assert.True(t, cfg.Has(test.key), "Has should return true after Set")

			// Verify Get method
			value, exists := test.getFunc(test.key)
			assert.True(t, exists, "Get should return true for existing key")
			assert.Equal(t, test.expected, value, "Get should return the correct value")

			// Verify Get for non-existent key
			_, exists = test.getFunc(test.key + ".nonexistent")
			assert.False(t, exists, "Get should return false for non-existent key")
		})
	}
}

func TestManager_GetTime(t *testing.T) {
	cfg := NewConfig()

	// Setup time value
	timeStr := "2023-05-24T15:00:00Z"
	err := cfg.Set("time.key", timeStr)
	assert.NoError(t, err)

	// Get time
	value, exists := cfg.GetTime("time.key")
	assert.True(t, exists)

	// Parse expected time
	expectedTime, _ := time.Parse(time.RFC3339, timeStr)
	assert.Equal(t, expectedTime, value)

	// Test non-existent key
	_, exists = cfg.GetTime("time.nonexistent")
	assert.False(t, exists)
}

func TestManager_GetSlice(t *testing.T) {
	cfg := NewConfig()

	// Test string slice
	stringSlice := []string{"a", "b", "c"}
	err := cfg.Set("slice.string", stringSlice)
	assert.NoError(t, err)

	sSlice, exists := cfg.GetStringSlice("slice.string")
	assert.True(t, exists)
	assert.Equal(t, stringSlice, sSlice)

	// Test int slice
	intSlice := []int{1, 2, 3}
	err = cfg.Set("slice.int", intSlice)
	assert.NoError(t, err)

	iSlice, exists := cfg.GetIntSlice("slice.int")
	assert.True(t, exists)
	assert.Equal(t, intSlice, iSlice)

	// Test generic slice
	genericSlice := []interface{}{1, "two", true}
	err = cfg.Set("slice.generic", genericSlice)
	assert.NoError(t, err)

	gSlice, exists := cfg.GetSlice("slice.generic")
	assert.True(t, exists)
	assert.Equal(t, genericSlice, gSlice)

	// Test invalid slice
	err = cfg.Set("not.a.slice", "string")
	assert.NoError(t, err)

	_, exists = cfg.GetSlice("not.a.slice")
	assert.False(t, exists)
}

func TestManager_GetMap(t *testing.T) {
	cfg := NewConfig()

	// Test direct map
	directMap := map[string]interface{}{
		"key1": "value1",
		"key2": 2,
		"key3": true,
	}
	err := cfg.Set("map.direct", directMap)
	assert.NoError(t, err)

	resultMap, exists := cfg.GetMap("map.direct")
	assert.True(t, exists)
	assert.Equal(t, directMap, resultMap)

	// Test map from subkeys
	err = cfg.Set("map.subkeys.key1", "value1")
	assert.NoError(t, err)
	err = cfg.Set("map.subkeys.key2", 2)
	assert.NoError(t, err)
	err = cfg.Set("map.subkeys.key3", true)
	assert.NoError(t, err)

	subkeyMap, exists := cfg.GetMap("map.subkeys")
	assert.True(t, exists)
	assert.Equal(t, map[string]interface{}{
		"key1": "value1",
		"key2": 2,
		"key3": true,
	}, subkeyMap)

	// Test specific map types
	stringMap := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}
	err = cfg.Set("map.string", stringMap)
	assert.NoError(t, err)

	sm, exists := cfg.GetStringMapString("map.string")
	assert.True(t, exists)
	assert.Equal(t, stringMap, sm)

	// Test string map string slice
	stringMapSlice := map[string][]string{
		"key1": {"a", "b", "c"},
		"key2": {"d", "e", "f"},
	}
	err = cfg.Set("map.stringslice", stringMapSlice)
	assert.NoError(t, err)

	sms, exists := cfg.GetStringMapStringSlice("map.stringslice")
	assert.True(t, exists)
	assert.Equal(t, stringMapSlice, sms)
}

func TestManager_SetDefault(t *testing.T) {
	cfg := NewConfig()

	// Set default value
	cfg.SetDefault("default.key", "default value")

	// Check value is returned
	value, exists := cfg.GetString("default.key")
	assert.True(t, exists)
	assert.Equal(t, "default value", value)

	// Override with Set
	err := cfg.Set("default.key", "overridden value")
	assert.NoError(t, err)

	// Check overridden value is returned
	value, exists = cfg.GetString("default.key")
	assert.True(t, exists)
	assert.Equal(t, "overridden value", value)
}

func TestManager_AllSettingsAndAllKeys(t *testing.T) {
	cfg := NewConfig()

	// Set some values
	err := cfg.Set("key1", "value1")
	assert.NoError(t, err)
	err = cfg.Set("section.key2", "value2")
	assert.NoError(t, err)
	err = cfg.Set("section.subsection.key3", "value3")
	assert.NoError(t, err)

	// Test AllSettings
	settings := cfg.AllSettings()
	assert.NotNil(t, settings)
	assert.Contains(t, settings, "key1")
	assert.Contains(t, settings, "section")

	// Test AllKeys
	keys := cfg.AllKeys()
	assert.NotNil(t, keys)
	assert.Contains(t, keys, "key1")
	assert.Contains(t, keys, "section.key2")
	assert.Contains(t, keys, "section.subsection.key3")
}

func TestManager_Unmarshal(t *testing.T) {
	cfg := NewConfig()

	// Define test structure
	type TestConfig struct {
		Name   string   `mapstructure:"name"`
		Count  int      `mapstructure:"count"`
		Items  []string `mapstructure:"items"`
		Nested struct {
			Enabled bool `mapstructure:"enabled"`
		} `mapstructure:"nested"`
	}

	// Set values
	err := cfg.Set("name", "test config")
	assert.NoError(t, err)
	err = cfg.Set("count", 42)
	assert.NoError(t, err)
	err = cfg.Set("items", []string{"item1", "item2", "item3"})
	assert.NoError(t, err)
	err = cfg.Set("nested.enabled", true)
	assert.NoError(t, err)

	// Unmarshal full config
	var fullConfig TestConfig
	err = cfg.Unmarshal(&fullConfig)
	assert.NoError(t, err)
	assert.Equal(t, "test config", fullConfig.Name)
	assert.Equal(t, 42, fullConfig.Count)
	assert.Equal(t, []string{"item1", "item2", "item3"}, fullConfig.Items)
	assert.True(t, fullConfig.Nested.Enabled)

	// Unmarshal nested section
	var nestedConfig struct {
		Enabled bool `mapstructure:"enabled"`
	}
	err = cfg.UnmarshalKey("nested", &nestedConfig)
	assert.NoError(t, err)
	assert.True(t, nestedConfig.Enabled)
}

func TestManager_ConfigFile(t *testing.T) {
	cfg := NewConfig()

	// Create a temporary YAML config file
	content := []byte(`
name: test config
count: 42
items:
  - item1
  - item2
  - item3
nested:
  enabled: true
`)

	// Write to temp file
	tmpfile, err := os.CreateTemp("", "config-*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	_, err = tmpfile.Write(content)
	require.NoError(t, err)
	err = tmpfile.Close()
	require.NoError(t, err)

	// Test SetConfigFile and ReadInConfig
	cfg.SetConfigFile(tmpfile.Name())
	err = cfg.ReadInConfig()
	assert.NoError(t, err)

	// Verify values were loaded
	name, exists := cfg.GetString("name")
	assert.True(t, exists)
	assert.Equal(t, "test config", name)

	count, exists := cfg.GetInt("count")
	assert.True(t, exists)
	assert.Equal(t, 42, count)

	// Test MergeInConfig
	// Create another config file to merge
	mergeContent := []byte(`
name: merged config
extra: extra value
`)

	mergeTmpfile, err := os.CreateTemp("", "merge-config-*.yaml")
	require.NoError(t, err)
	defer os.Remove(mergeTmpfile.Name())

	_, err = mergeTmpfile.Write(mergeContent)
	require.NoError(t, err)
	err = mergeTmpfile.Close()
	require.NoError(t, err)

	// Merge the config
	cfg.SetConfigFile(mergeTmpfile.Name())
	err = cfg.MergeInConfig()
	assert.NoError(t, err)

	// Verify merged values
	name, exists = cfg.GetString("name")
	assert.True(t, exists)
	assert.Equal(t, "merged config", name)

	extra, exists := cfg.GetString("extra")
	assert.True(t, exists)
	assert.Equal(t, "extra value", extra)

	// Original values should still be there
	count, exists = cfg.GetInt("count")
	assert.True(t, exists)
	assert.Equal(t, 42, count)
}

func TestManager_ConfigType(t *testing.T) {
	cfg := NewConfig()

	// Test SetConfigType and read from buffer
	cfg.SetConfigType("json")

	jsonConfig := []byte(`{
		"name": "json config",
		"value": 123
	}`)

	err := cfg.MergeConfig(bytes.NewBuffer(jsonConfig))
	assert.NoError(t, err)

	name, exists := cfg.GetString("name")
	assert.True(t, exists)
	assert.Equal(t, "json config", name)

	value, exists := cfg.GetInt("value")
	assert.True(t, exists)
	assert.Equal(t, 123, value)
}

func TestManager_SetConfigName(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "config-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create test config file
	configPath := tmpDir + "/config.json"
	err = os.WriteFile(configPath, []byte(`{"test": "value"}`), 0644)
	require.NoError(t, err)

	cfg := NewConfig()
	cfg.SetConfigName("config")
	cfg.AddConfigPath(tmpDir)

	err = cfg.ReadInConfig()
	assert.NoError(t, err)

	value, exists := cfg.GetString("test")
	assert.True(t, exists)
	assert.Equal(t, "value", value)
}

func TestManager_WriteConfig(t *testing.T) {
	cfg := NewConfig()

	// Set some values
	err := cfg.Set("test", "value")
	assert.NoError(t, err)

	// Create temp file path
	tmpfile, err := os.CreateTemp("", "write-config-*.json")
	require.NoError(t, err)
	err = tmpfile.Close()
	require.NoError(t, err)
	os.Remove(tmpfile.Name()) // Remove so WriteConfigAs doesn't fail

	// Test WriteConfigAs
	cfg.SetConfigType("json")
	err = cfg.WriteConfigAs(tmpfile.Name())
	assert.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	// Read the file and check content
	content, err := os.ReadFile(tmpfile.Name())
	assert.NoError(t, err)
	assert.Contains(t, string(content), "test")
	assert.Contains(t, string(content), "value")

	// Test SafeWriteConfigAs (should fail since file exists)
	err = cfg.SafeWriteConfigAs(tmpfile.Name())
	assert.Error(t, err)

	// Test with a new name
	newTmpfile, err := os.CreateTemp("", "safe-write-config-*.json")
	require.NoError(t, err)
	newName := newTmpfile.Name()
	err = newTmpfile.Close()
	require.NoError(t, err)
	os.Remove(newName) // Remove so SafeWriteConfigAs doesn't fail

	err = cfg.SafeWriteConfigAs(newName)
	assert.NoError(t, err)
	defer os.Remove(newName)
}

func TestManager_EnvSupport(t *testing.T) {
	// Setup - direct key format since we can't use SetEnvKeyReplacer
	os.Setenv("TESTCONFIG_TEST_KEY", "env value")
	defer os.Unsetenv("TESTCONFIG_TEST_KEY")

	cfg := NewConfig()
	cfg.SetEnvPrefix("TESTCONFIG")
	// Viper sẽ tự động chuyển đổi dấu chấm thành dấu gạch dưới
	cfg.AutomaticEnv()

	// Bind environment variable directly
	err := cfg.BindEnv("test.key", "TESTCONFIG_TEST_KEY")
	assert.NoError(t, err)

	// Test bound env var
	value, exists := cfg.GetString("test.key")
	assert.True(t, exists, "Environment variable should be accessible through config")
	assert.Equal(t, "env value", value)

	// Test another binding
	os.Setenv("CUSTOM_ENV", "custom value")
	defer os.Unsetenv("CUSTOM_ENV")

	err = cfg.BindEnv("bound.key", "CUSTOM_ENV")
	assert.NoError(t, err)

	value, exists = cfg.GetString("bound.key")
	assert.True(t, exists, "Bound environment variable should be accessible")
	assert.Equal(t, "custom value", value)
}

func TestManager_WatchConfig(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping in CI environment")
	}

	// Create temp config file
	tmpfile, err := os.CreateTemp("", "watch-config-*.json")
	require.NoError(t, err)

	_, err = tmpfile.Write([]byte(`{"test": "original"}`))
	require.NoError(t, err)
	err = tmpfile.Close()
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	cfg := NewConfig()
	cfg.SetConfigFile(tmpfile.Name())
	err = cfg.ReadInConfig()
	assert.NoError(t, err)

	// Setup change detection
	changeDetected := make(chan bool)
	cfg.OnConfigChange(func(event fsnotify.Event) {
		changeDetected <- true
	})
	cfg.WatchConfig()

	// Update file
	time.Sleep(100 * time.Millisecond) // Wait a bit for watcher to be ready
	err = os.WriteFile(tmpfile.Name(), []byte(`{"test": "updated"}`), 0644)
	assert.NoError(t, err)

	// Wait for change event or timeout
	select {
	case <-changeDetected:
		// Config change was detected
		value, exists := cfg.GetString("test")
		assert.True(t, exists)
		assert.Equal(t, "updated", value)
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout waiting for config change event")
	}
}

// Test thêm các trường hợp còn thiếu
func TestManager_ExtendedMapTests(t *testing.T) {
	cfg := NewConfig()

	// Test cho GetMap với nil value và key không tồn tại
	// Kiểm tra trường hợp key không tồn tại
	resultMap, exists := cfg.GetMap("nonexistent.map")
	assert.False(t, exists)
	assert.Nil(t, resultMap)

	// Tạo một key với giá trị không phải map
	err := cfg.Set("not.a.map", "string value")
	assert.NoError(t, err)
	resultMap, exists = cfg.GetMap("not.a.map")
	assert.False(t, exists)
	assert.Nil(t, resultMap)

	// Test GetStringMap
	// Setup a string map
	stringMap := map[string]interface{}{
		"key1": "value1",
		"key2": 2,
	}
	err = cfg.Set("string.map", stringMap)
	assert.NoError(t, err)

	// Test với key tồn tại
	sMap, exists := cfg.GetStringMap("string.map")
	assert.True(t, exists)
	assert.Equal(t, stringMap, sMap)

	// Test với key không tồn tại
	sMap, exists = cfg.GetStringMap("nonexistent.string.map")
	assert.False(t, exists)
	assert.Nil(t, sMap)
}

func TestManager_ExtendedFileOps(t *testing.T) {
	cfg := NewConfig()

	// Setup a config file for WriteConfig and SafeWriteConfig
	tmpfile, err := os.CreateTemp("", "config-*.yaml")
	assert.NoError(t, err)
	tmpfileName := tmpfile.Name()
	err = tmpfile.Close()
	assert.NoError(t, err)

	// Set some values
	err = cfg.Set("test", "value")
	assert.NoError(t, err)

	// Set config file
	cfg.SetConfigFile(tmpfileName)

	// Test WriteConfig
	err = cfg.WriteConfig()
	assert.NoError(t, err)

	// Check content
	content, err := os.ReadFile(tmpfileName)
	assert.NoError(t, err)
	assert.Contains(t, string(content), "test")

	// Test SafeWriteConfig (should fail since file exists)
	err = cfg.SafeWriteConfig()
	assert.Error(t, err, "SafeWriteConfig should fail if file exists")

	// Clean up
	os.Remove(tmpfileName)
}

func TestServiceProvider_BootComplete(t *testing.T) {
	provider := NewServiceProvider()

	// Test with nil app
	assert.NotPanics(t, func() {
		provider.Boot(nil)
	}, "Boot should not panic with nil app")
}

func TestManager_MoreSliceTests(t *testing.T) {
	cfg := NewConfig()

	// Test cho các trường hợp còn thiếu
	// GetSlice với key không tồn tại
	slice, exists := cfg.GetSlice("nonexistent.slice")
	assert.False(t, exists)
	assert.Nil(t, slice)

	// GetStringSlice với key không tồn tại
	strSlice, exists := cfg.GetStringSlice("nonexistent.stringslice")
	assert.False(t, exists)
	assert.Nil(t, strSlice)

	// GetIntSlice với key không tồn tại
	intSlice, exists := cfg.GetIntSlice("nonexistent.intslice")
	assert.False(t, exists)
	assert.Nil(t, intSlice)

	// GetStringMapString với key không tồn tại
	strMapStr, exists := cfg.GetStringMapString("nonexistent.stringmapstring")
	assert.False(t, exists)
	assert.Nil(t, strMapStr)

	// GetStringMapStringSlice với key không tồn tại
	strMapStrSlice, exists := cfg.GetStringMapStringSlice("nonexistent.stringmapstringslice")
	assert.False(t, exists)
	assert.Nil(t, strMapStrSlice)
}

func Test_getMapFromSubKeys_NoSubKey(t *testing.T) {
	cfg := NewConfig()
	m := cfg.(*manager)
	result, ok := m.getMapFromSubKeys("foo")
	assert.False(t, ok)
	assert.Nil(t, result)
}

func Test_getMapFromSubKeys_PartialPrefixButNotSubKey(t *testing.T) {
	cfg := NewConfig()
	_ = cfg.Set("foobar.key1", "value1")
	m := cfg.(*manager)
	result, ok := m.getMapFromSubKeys("foo")
	assert.False(t, ok)
	assert.Nil(t, result)
}

func Test_getMapFromSubKeys_SubKeyExists(t *testing.T) {
	cfg := NewConfig()
	_ = cfg.Set("foo.key1", "v1")
	_ = cfg.Set("foo.key2", 2)
	_ = cfg.Set("foo.key3", true)
	m := cfg.(*manager)
	result, ok := m.getMapFromSubKeys("foo")
	assert.True(t, ok)
	assert.Equal(t, map[string]interface{}{
		"key1": "v1",
		"key2": 2,
		"key3": true,
	}, result)
}
