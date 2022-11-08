package pubsub

const NotImplemented = PubsubError("pubsub: not implemented")                      // nolint:errname
const WrongDataType = PubsubError("pubsub: wrong data type")                       // nolint:errname
const ListenWhileConsuming = PubsubError("pubsub: listen while consuming")         // nolint:errname
const ConsumeWhileListening = PubsubError("pubsub: consume while listening")       // nolint:errname
const TopicNotExist = PubsubError("pubsub: topic not exist")                       // nolint:errname
const ListenNotPointer = PubsubError("pubsub: listen model is not pointer")        // nolint:errname
const ListenNullCallback = PubsubError("pubsub: listen callback is null")          // nolint:errname
const TypeCastingFailed = PubsubError("pubsub: type casting failed")               // nolint:errname
const InputNotString = PubsubError("pubsub: input model is not pointer of string") // nolint:errname
const NullDelivery = PubsubError("pubsub: not having any delivery")                // nolint:errname

type PubsubError string

func (e PubsubError) Error() string { return string(e) }

func (PubsubError) PubsubError() {}
