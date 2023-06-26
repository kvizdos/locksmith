# What does this directory DO?

The Locksmith `structs` package allows developers to easily encrypt and decrypt data that matters, all by specifying a structure tag.

The encryption will encode all structure values into their JSON or BSON (serializer) tag values, dependent on which is set and chosen.

During encryption, if a serializer tag is missing or set to "-", it will skip the item from being included in the final map
