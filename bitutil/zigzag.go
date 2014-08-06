package bitutil

// DecodeZigZag decodes a ZigZag encoded integer and returns it.
//
// Examples of ZigZag encoded values on the left and decoded values on the
// right:
//
//     0 => 0
//    -1 => 1
//     1 => 2
//    -2 => 3
//     2 => 4
//    -3 => 5
//     3 => 6
//
// ref: https://developers.google.com/protocol-buffers/docs/encoding
func DecodeZigZag(x int32) int32 {
	return x>>1 ^ -(x & 1)
}
