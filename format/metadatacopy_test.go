package format

import (
	"testing"

	"github.com/trivago/gollum/core"
	"github.com/trivago/tgo/ttesting"
)

func TestMetadataCopyReplace(t *testing.T) {
	expect := ttesting.NewExpect(t)

	config := core.NewPluginConfig("", "format.MetadataCopy")
	config.Override("Key", "foo")

	plugin, err := core.NewPluginWithConfig(config)
	expect.NoError(err)

	formatter, casted := plugin.(*MetadataCopy)
	expect.True(casted)

	msg := core.NewMessage(nil, []byte("test"), core.Metadata{"foo": []byte("foo")}, core.InvalidStreamID)

	err = formatter.ApplyFormatter(msg)
	expect.NoError(err)

	expect.Equal("foo", msg.String())
}

func TestMetadataCopyAddKey(t *testing.T) {
	expect := ttesting.NewExpect(t)

	config := core.NewPluginConfig("", "format.MetadataCopy")
	config.Override("ApplyTo", "foo")

	plugin, err := core.NewPluginWithConfig(config)
	expect.NoError(err)

	formatter, casted := plugin.(*MetadataCopy)
	expect.True(casted)

	msg := core.NewMessage(nil, []byte("test"), nil, core.InvalidStreamID)

	err = formatter.ApplyFormatter(msg)
	expect.NoError(err)

	val, err := msg.GetMetadata().String("foo")
	expect.NoError(err)
	expect.Equal("test", msg.String())
	expect.Equal("test", val)
}

func TestMetadataCopyAppend(t *testing.T) {
	expect := ttesting.NewExpect(t)

	config := core.NewPluginConfig("", "format.MetadataCopy")
	config.Override("Key", "foo")
	config.Override("Mode", "append")
	config.Override("Separator", " ")

	plugin, err := core.NewPluginWithConfig(config)
	expect.NoError(err)

	formatter, casted := plugin.(*MetadataCopy)
	expect.True(casted)

	msg := core.NewMessage(nil, []byte("test"), core.Metadata{"foo": []byte("foo")}, core.InvalidStreamID)

	err = formatter.ApplyFormatter(msg)
	expect.NoError(err)

	expect.Equal("test foo", msg.String())
}

func TestMetadataCopyPrepend(t *testing.T) {
	expect := ttesting.NewExpect(t)

	config := core.NewPluginConfig("", "format.MetadataCopy")
	config.Override("Key", "foo")
	config.Override("Mode", "prepend")

	plugin, err := core.NewPluginWithConfig(config)
	expect.NoError(err)

	formatter, casted := plugin.(*MetadataCopy)
	expect.True(casted)

	msg := core.NewMessage(nil, []byte("test"), core.Metadata{"foo": []byte("foo")}, core.InvalidStreamID)

	err = formatter.ApplyFormatter(msg)
	expect.NoError(err)

	expect.Equal("footest", msg.String())
}

func TestMetadataCopyDeprecated(t *testing.T) {
	expect := ttesting.NewExpect(t)

	config := core.NewPluginConfig("", "format.MetadataCopy")
	config.Override("CopyToKeys", []string{"foo", "bar"})

	plugin, err := core.NewPluginWithConfig(config)
	expect.NoError(err)

	formatter, casted := plugin.(*MetadataCopy)
	expect.True(casted)

	msg := core.NewMessage(nil, []byte("test"), nil, core.InvalidStreamID)

	err = formatter.ApplyFormatter(msg)
	expect.NoError(err)

	foo, err := msg.GetMetadata().String("foo")
	expect.NoError(err)
	bar, err := msg.GetMetadata().String("bar")
	expect.NoError(err)

	expect.Equal("test", msg.String())
	expect.Equal("test", foo)
	expect.Equal("test", bar)
}

func TestMetadataCopyApplyToHandlingDeprecated(t *testing.T) {
	expect := ttesting.NewExpect(t)

	config := core.NewPluginConfig("", "format.MetadataCopy")
	config.Override("ApplyTo", "foo")
	config.Override("CopyToKeys", []string{"bar"})

	plugin, err := core.NewPluginWithConfig(config)
	expect.NoError(err)

	formatter, casted := plugin.(*MetadataCopy)
	expect.True(casted)

	msg := core.NewMessage(nil, []byte("payload"), nil, core.InvalidStreamID)
	msg.GetMetadata().Set("foo", []byte("meta"))

	err = formatter.ApplyFormatter(msg)
	expect.NoError(err)

	foo, err := msg.GetMetadata().String("foo")
	expect.NoError(err)
	bar, err := msg.GetMetadata().String("bar")
	expect.NoError(err)

	expect.Equal("payload", msg.String())
	expect.Equal("meta", foo)
	expect.Equal("meta", bar)
}

func TestMetadataCopyMetadataIntegrity(t *testing.T) {
	expect := ttesting.NewExpect(t)

	config := core.NewPluginConfig("", "format.MetadataCopy")
	config.Override("ApplyTo", "foo")

	plugin, err := core.NewPluginWithConfig(config)
	expect.NoError(err)

	formatter, casted := plugin.(*MetadataCopy)
	expect.True(casted)

	msg := core.NewMessage(nil, []byte("payload"), nil, core.InvalidStreamID)

	err = formatter.ApplyFormatter(msg)
	expect.NoError(err)

	foo, err := msg.GetMetadata().String("foo")
	expect.NoError(err)
	expect.Equal("payload", msg.String())
	expect.Equal("payload", foo)

	msg.StorePayload([]byte("xxx"))

	foo, err = msg.GetMetadata().String("foo")
	expect.NoError(err)
	expect.Equal("xxx", msg.String())
	expect.Equal("payload", foo)
}

func TestMetadataCopyPayloadIntegrity(t *testing.T) {
	expect := ttesting.NewExpect(t)

	config := core.NewPluginConfig("", "format.MetadataCopy")
	config.Override("Key", "foo")

	plugin, err := core.NewPluginWithConfig(config)
	expect.NoError(err)

	formatter, casted := plugin.(*MetadataCopy)
	expect.True(casted)

	msg := core.NewMessage(nil, []byte{}, nil, core.InvalidStreamID)
	msg.GetMetadata().Set("foo", []byte("metadata"))

	err = formatter.ApplyFormatter(msg)
	expect.NoError(err)

	foo, err := msg.GetMetadata().Bytes("foo")
	expect.NoError(err)
	expect.Equal("metadata", msg.String())
	expect.Equal("metadata", string(foo))

	msg.GetMetadata().Set("foo", []byte("xxx"))

	foo, err = msg.GetMetadata().Bytes("foo")
	expect.NoError(err)
	expect.Equal("metadata", msg.String())
	expect.Equal("xxx", string(foo))
}
