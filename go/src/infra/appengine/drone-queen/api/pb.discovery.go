// Code generated by cproto. DO NOT EDIT.

package api

import "go.chromium.org/luci/grpc/discovery"

import "google.golang.org/protobuf/types/descriptorpb"

func init() {
	discovery.RegisterDescriptorSetCompressed(
		[]string{
			"drone_queen.Drone", "drone_queen.InventoryProvider", "drone_queen.Inspect",
		},
		[]byte{31, 139,
			8, 0, 0, 0, 0, 0, 0, 255, 164, 90, 77, 108, 27, 73,
			118, 102, 255, 72, 150, 74, 254, 83, 217, 227, 209, 208, 30, 251,
			153, 30, 237, 72, 179, 84, 83, 146, 127, 38, 178, 119, 50, 75,
			137, 148, 69, 175, 68, 42, 252, 177, 199, 246, 238, 218, 37, 118,
			145, 236, 153, 102, 55, 167, 171, 40, 89, 246, 122, 146, 108, 46,
			89, 32, 72, 128, 28, 130, 0, 65, 128, 252, 0, 9, 176, 57,
			231, 148, 92, 114, 207, 45, 64, 128, 156, 114, 218, 251, 2, 185,
			228, 24, 188, 170, 238, 38, 37, 217, 227, 157, 141, 48, 198, 244,
			171, 159, 247, 243, 213, 171, 247, 94, 85, 145, 252, 71, 129, 92,
			235, 134, 97, 215, 231, 133, 65, 20, 202, 112, 111, 216, 41, 72,
			175, 207, 133, 100, 253, 129, 163, 154, 232, 57, 61, 192, 73, 6,
			228, 238, 145, 233, 102, 50, 134, 206, 145, 83, 130, 183, 195, 192,
			21, 115, 6, 24, 11, 86, 61, 33, 233, 69, 50, 17, 176, 32,
			20, 115, 38, 24, 11, 19, 117, 77, 172, 127, 67, 46, 180, 195,
			190, 115, 140, 231, 250, 217, 148, 227, 46, 54, 237, 26, 79, 190,
			223, 245, 100, 111, 184, 231, 180, 195, 126, 161, 27, 250, 44, 232,
			142, 84, 28, 200, 195, 1, 23, 35, 77, 255, 215, 48, 254, 206,
			180, 238, 239, 174, 255, 210, 188, 122, 95, 115, 222, 141, 199, 58,
			143, 184, 239, 255, 40, 8, 15, 130, 38, 206, 121, 240, 95, 75,
			228, 20, 157, 184, 154, 249, 133, 97, 144, 127, 63, 77, 140, 211,
			212, 186, 154, 161, 171, 255, 118, 26, 212, 140, 118, 232, 195, 250,
			176, 211, 225, 145, 128, 37, 208, 188, 62, 22, 224, 50, 201, 192,
			11, 36, 143, 218, 61, 22, 116, 57, 116, 194, 168, 207, 36, 129,
			141, 112, 112, 24, 121, 221, 158, 132, 213, 229, 229, 223, 137, 39,
			64, 37, 104, 59, 0, 69, 223, 7, 213, 39, 32, 226, 130, 71,
			251, 220, 117, 8, 244, 164, 28, 136, 187, 133, 130, 203, 247, 185,
			31, 14, 120, 36, 18, 48, 208, 210, 65, 172, 196, 210, 158, 86,
			162, 64, 8, 212, 185, 235, 9, 25, 121, 123, 67, 233, 133, 1,
			176, 192, 133, 161, 224, 224, 5, 32, 194, 97, 212, 230, 170, 101,
			207, 11, 88, 116, 168, 244, 18, 121, 56, 240, 100, 15, 194, 72,
			253, 63, 28, 74, 2, 253, 208, 245, 58, 94, 155, 33, 135, 60,
			176, 136, 195, 128, 71, 125, 79, 74, 238, 194, 32, 10, 247, 61,
			151, 187, 32, 123, 76, 130, 236, 161, 117, 190, 31, 30, 120, 65,
			23, 112, 41, 61, 156, 36, 112, 18, 129, 62, 151, 119, 9, 1,
			252, 251, 228, 152, 98, 2, 194, 78, 162, 81, 59, 116, 57, 244,
			135, 66, 66, 196, 37, 243, 2, 197, 149, 237, 133, 251, 216, 21,
			35, 70, 32, 8, 165, 215, 230, 121, 144, 61, 79, 128, 239, 9,
			137, 28, 198, 37, 6, 238, 49, 117, 92, 79, 180, 125, 230, 245,
			121, 228, 188, 77, 9, 47, 24, 199, 34, 81, 98, 16, 133, 238,
			176, 205, 71, 122, 144, 145, 34, 255, 47, 61, 8, 196, 214, 185,
			97, 123, 216, 231, 129, 100, 201, 34, 21, 194, 8, 66, 217, 227,
			17, 244, 153, 228, 145, 199, 124, 49, 130, 90, 45, 144, 236, 113,
			2, 227, 218, 167, 70, 85, 185, 167, 102, 34, 227, 128, 245, 57,
			42, 52, 238, 91, 65, 56, 234, 83, 184, 123, 82, 160, 69, 129,
			102, 21, 70, 2, 250, 236, 16, 246, 56, 122, 138, 11, 50, 4,
			30, 184, 97, 36, 56, 58, 197, 32, 10, 251, 161, 228, 160, 49,
			145, 2, 92, 30, 121, 251, 220, 133, 78, 20, 246, 137, 70, 65,
			132, 29, 121, 128, 110, 18, 123, 16, 136, 1, 111, 163, 7, 193,
			32, 242, 208, 177, 34, 244, 157, 64, 123, 145, 16, 74, 119, 2,
			205, 173, 74, 3, 26, 181, 205, 230, 163, 98, 189, 12, 149, 6,
			236, 214, 107, 15, 43, 165, 114, 9, 214, 31, 67, 115, 171, 12,
			27, 181, 221, 199, 245, 202, 253, 173, 38, 108, 213, 182, 75, 229,
			122, 3, 138, 213, 18, 108, 212, 170, 205, 122, 101, 189, 213, 172,
			213, 27, 4, 114, 197, 6, 84, 26, 57, 213, 83, 172, 62, 134,
			242, 23, 187, 245, 114, 163, 1, 181, 58, 84, 118, 118, 183, 43,
			229, 18, 60, 42, 214, 235, 197, 106, 179, 82, 110, 228, 161, 82,
			221, 216, 110, 149, 42, 213, 251, 121, 88, 111, 53, 161, 90, 107,
			18, 216, 174, 236, 84, 154, 229, 18, 52, 107, 121, 37, 246, 228,
			60, 168, 109, 194, 78, 185, 190, 177, 85, 172, 54, 139, 235, 149,
			237, 74, 243, 177, 18, 184, 89, 105, 86, 81, 216, 102, 173, 78,
			160, 8, 187, 197, 122, 179, 178, 209, 218, 46, 214, 97, 183, 85,
			223, 173, 53, 202, 128, 150, 149, 42, 141, 141, 237, 98, 101, 167,
			92, 114, 160, 82, 133, 106, 13, 202, 15, 203, 213, 38, 52, 182,
			138, 219, 219, 71, 13, 37, 80, 123, 84, 45, 215, 81, 251, 113,
			51, 97, 189, 12, 219, 149, 226, 250, 118, 25, 69, 41, 59, 75,
			149, 122, 121, 163, 137, 6, 141, 190, 54, 42, 165, 114, 181, 89,
			220, 206, 19, 104, 236, 150, 55, 42, 197, 237, 60, 148, 191, 40,
			239, 236, 110, 23, 235, 143, 243, 49, 211, 70, 249, 247, 90, 229,
			106, 179, 82, 220, 134, 82, 113, 167, 120, 191, 220, 128, 133, 119,
			161, 178, 91, 175, 109, 180, 234, 229, 29, 212, 186, 182, 9, 141,
			214, 122, 163, 89, 105, 182, 154, 101, 184, 95, 171, 149, 20, 216,
			141, 114, 253, 97, 101, 163, 220, 184, 7, 219, 181, 134, 2, 172,
			213, 40, 231, 9, 148, 138, 205, 162, 18, 189, 91, 175, 109, 86,
			154, 141, 123, 248, 189, 222, 106, 84, 20, 112, 149, 106, 179, 92,
			175, 183, 118, 155, 149, 90, 117, 17, 182, 106, 143, 202, 15, 203,
			117, 216, 40, 182, 26, 229, 146, 66, 184, 86, 69, 107, 209, 87,
			202, 181, 250, 99, 100, 139, 56, 168, 21, 200, 195, 163, 173, 114,
			115, 171, 92, 71, 80, 21, 90, 69, 132, 161, 209, 172, 87, 54,
			154, 227, 195, 106, 117, 104, 214, 234, 77, 50, 102, 39, 84, 203,
			247, 183, 43, 247, 203, 213, 141, 50, 118, 215, 144, 205, 163, 74,
			163, 188, 8, 197, 122, 165, 129, 3, 42, 74, 48, 60, 42, 62,
			134, 90, 75, 89, 141, 11, 213, 106, 148, 137, 254, 30, 115, 221,
			188, 90, 79, 168, 108, 66, 177, 244, 176, 130, 154, 199, 163, 119,
			107, 141, 70, 37, 118, 23, 5, 219, 198, 86, 140, 185, 67, 200,
			20, 49, 76, 106, 65, 102, 14, 191, 166, 168, 149, 203, 220, 35,
			211, 196, 156, 154, 215, 159, 186, 241, 70, 230, 154, 106, 188, 166,
			63, 117, 227, 71, 153, 117, 213, 56, 163, 63, 117, 227, 124, 38,
			175, 26, 13, 253, 169, 27, 191, 151, 41, 168, 198, 248, 83, 55,
			126, 156, 201, 169, 70, 162, 63, 117, 227, 66, 230, 186, 106, 252,
			72, 127, 254, 207, 101, 98, 218, 25, 58, 241, 13, 102, 190, 236,
			175, 46, 67, 17, 210, 148, 171, 226, 35, 23, 60, 144, 2, 24,
			12, 66, 47, 144, 42, 170, 121, 125, 204, 50, 46, 31, 240, 192,
			229, 129, 138, 138, 44, 56, 212, 237, 47, 195, 64, 5, 19, 63,
			108, 51, 159, 64, 155, 249, 60, 112, 89, 148, 7, 30, 96, 240,
			119, 129, 33, 175, 118, 56, 212, 243, 226, 162, 64, 133, 210, 78,
			196, 218, 163, 132, 145, 116, 96, 62, 192, 10, 65, 209, 152, 48,
			67, 95, 199, 68, 104, 246, 120, 204, 200, 195, 76, 234, 51, 233,
			237, 115, 140, 105, 44, 0, 62, 8, 219, 61, 96, 18, 90, 205,
			13, 232, 123, 110, 160, 2, 122, 24, 16, 120, 192, 130, 33, 102,
			129, 149, 60, 172, 172, 125, 186, 156, 79, 226, 244, 32, 10, 125,
			62, 144, 94, 27, 238, 71, 188, 27, 70, 30, 11, 82, 237, 225,
			160, 231, 181, 123, 192, 95, 72, 142, 58, 169, 248, 252, 134, 81,
			123, 172, 253, 213, 1, 139, 112, 68, 8, 135, 156, 69, 16, 6,
			28, 227, 31, 102, 252, 190, 23, 12, 37, 87, 233, 18, 238, 44,
			167, 246, 249, 97, 208, 117, 96, 155, 179, 193, 200, 228, 136, 67,
			78, 244, 57, 139, 184, 155, 3, 17, 234, 252, 27, 132, 224, 115,
			54, 32, 241, 48, 144, 108, 207, 231, 104, 121, 192, 57, 226, 218,
			9, 35, 93, 137, 12, 48, 181, 234, 124, 62, 20, 152, 148, 24,
			60, 93, 189, 181, 212, 11, 135, 17, 248, 94, 192, 89, 68, 64,
			113, 255, 201, 194, 183, 215, 28, 184, 158, 5, 53, 114, 81, 5,
			241, 30, 135, 72, 21, 57, 158, 80, 41, 1, 150, 151, 151, 87,
			150, 212, 127, 205, 229, 229, 187, 234, 191, 39, 104, 250, 218, 218,
			218, 218, 210, 202, 234, 210, 205, 149, 230, 234, 205, 187, 183, 215,
			238, 222, 94, 115, 214, 146, 191, 39, 14, 172, 31, 18, 92, 72,
			25, 121, 109, 137, 10, 202, 216, 68, 197, 61, 15, 7, 28, 120,
			32, 134, 17, 215, 173, 7, 28, 218, 136, 114, 24, 236, 243, 72,
			234, 245, 213, 57, 9, 158, 214, 55, 55, 8, 220, 188, 121, 115,
			109, 100, 203, 193, 193, 129, 227, 113, 217, 113, 194, 168, 91, 136,
			58, 109, 252, 135, 35, 28, 249, 66, 46, 98, 193, 198, 1, 37,
			7, 93, 129, 70, 221, 128, 242, 11, 214, 31, 248, 92, 16, 146,
			124, 194, 202, 93, 216, 8, 251, 131, 161, 228, 99, 123, 65, 9,
			220, 173, 53, 42, 95, 192, 115, 68, 102, 97, 241, 185, 19, 87,
			60, 163, 65, 105, 229, 121, 79, 247, 140, 106, 102, 193, 229, 179,
			120, 129, 23, 212, 244, 106, 107, 123, 123, 113, 241, 141, 227, 148,
			191, 47, 44, 47, 222, 27, 211, 105, 245, 93, 58, 117, 185, 68,
			46, 97, 199, 101, 135, 99, 186, 9, 25, 13, 219, 82, 9, 216,
			103, 62, 200, 253, 88, 226, 145, 225, 223, 147, 251, 121, 80, 10,
			221, 251, 109, 77, 218, 119, 228, 62, 82, 223, 102, 145, 30, 52,
			20, 188, 13, 159, 192, 202, 242, 242, 81, 11, 111, 190, 213, 194,
			71, 94, 112, 115, 21, 158, 223, 231, 178, 113, 40, 36, 239, 99,
			119, 81, 108, 122, 62, 111, 30, 93, 136, 205, 202, 118, 185, 89,
			217, 41, 67, 71, 198, 106, 188, 109, 206, 247, 58, 50, 209, 180,
			85, 169, 54, 239, 220, 2, 233, 181, 191, 18, 240, 25, 44, 44,
			44, 232, 150, 197, 142, 116, 220, 131, 45, 175, 219, 43, 49, 169,
			102, 45, 194, 15, 126, 0, 55, 87, 23, 225, 103, 160, 250, 182,
			195, 131, 164, 43, 193, 173, 80, 128, 34, 234, 235, 134, 7, 66,
			177, 196, 205, 178, 178, 188, 60, 22, 195, 132, 147, 14, 208, 81,
			106, 229, 206, 201, 109, 148, 114, 195, 233, 43, 119, 110, 221, 186,
			245, 233, 205, 59, 203, 163, 176, 177, 199, 59, 97, 196, 161, 21,
			120, 47, 18, 46, 107, 159, 46, 31, 231, 226, 252, 118, 139, 185,
			160, 237, 135, 133, 5, 13, 74, 65, 45, 22, 254, 45, 194, 210,
			184, 58, 239, 240, 96, 228, 131, 112, 37, 124, 230, 199, 248, 40,
			7, 88, 60, 226, 0, 183, 222, 234, 0, 15, 216, 62, 131, 231,
			122, 33, 157, 246, 48, 138, 120, 32, 113, 200, 142, 231, 251, 158,
			24, 115, 0, 140, 166, 208, 87, 173, 240, 25, 188, 125, 194, 183,
			184, 57, 124, 54, 106, 117, 2, 126, 176, 62, 244, 124, 151, 71,
			11, 139, 104, 88, 35, 70, 40, 22, 161, 129, 89, 212, 188, 240,
			15, 199, 84, 181, 237, 94, 32, 209, 242, 120, 164, 54, 61, 54,
			91, 33, 176, 232, 236, 33, 103, 165, 203, 8, 131, 219, 111, 197,
			32, 182, 34, 201, 190, 176, 123, 40, 123, 186, 186, 62, 2, 255,
			184, 250, 11, 139, 199, 215, 230, 62, 151, 27, 35, 52, 22, 22,
			85, 4, 124, 208, 168, 85, 97, 135, 13, 6, 94, 208, 37, 4,
			42, 129, 110, 209, 71, 217, 188, 74, 142, 99, 56, 29, 14, 84,
			2, 56, 146, 206, 117, 64, 141, 51, 41, 81, 97, 249, 59, 69,
			101, 45, 10, 51, 58, 195, 100, 158, 215, 108, 116, 43, 10, 203,
			189, 194, 108, 250, 122, 233, 85, 63, 12, 100, 239, 245, 210, 43,
			151, 29, 190, 110, 190, 194, 148, 246, 250, 238, 171, 190, 23, 188,
			190, 251, 74, 240, 246, 235, 167, 206, 43, 44, 34, 208, 145, 95,
			255, 228, 73, 142, 192, 65, 143, 71, 28, 244, 108, 100, 196, 252,
			3, 118, 40, 128, 191, 192, 186, 6, 79, 64, 58, 67, 118, 48,
			55, 186, 94, 215, 147, 2, 83, 189, 207, 33, 150, 148, 7, 37,
			42, 79, 64, 11, 203, 131, 146, 150, 87, 41, 72, 137, 84, 217,
			250, 37, 143, 194, 165, 1, 115, 93, 125, 166, 146, 7, 97, 194,
			141, 179, 118, 79, 87, 42, 73, 117, 131, 85, 81, 188, 209, 242,
			113, 93, 129, 233, 173, 27, 194, 112, 160, 146, 103, 50, 117, 193,
			115, 184, 19, 55, 174, 188, 185, 6, 90, 204, 19, 37, 63, 28,
			104, 206, 90, 82, 238, 73, 14, 196, 176, 211, 241, 94, 96, 149,
			134, 135, 123, 174, 106, 22, 229, 7, 170, 62, 91, 200, 181, 154,
			27, 185, 197, 123, 71, 90, 137, 46, 163, 190, 30, 122, 17, 119,
			29, 40, 130, 186, 115, 184, 169, 157, 65, 168, 131, 170, 247, 146,
			71, 32, 122, 225, 208, 119, 19, 40, 135, 130, 171, 26, 107, 129,
			137, 84, 154, 11, 123, 135, 4, 213, 88, 196, 5, 8, 240, 104,
			24, 232, 68, 127, 210, 149, 16, 72, 118, 68, 212, 128, 69, 98,
			36, 102, 143, 19, 80, 149, 14, 230, 253, 118, 155, 15, 36, 236,
			133, 178, 167, 100, 226, 92, 125, 146, 78, 108, 16, 39, 244, 192,
			98, 48, 236, 116, 4, 151, 170, 136, 217, 12, 35, 224, 122, 175,
			229, 33, 183, 186, 188, 242, 41, 198, 204, 149, 219, 205, 229, 149,
			187, 55, 151, 239, 174, 220, 118, 150, 87, 158, 228, 98, 239, 22,
			160, 232, 52, 232, 14, 152, 144, 4, 212, 72, 37, 63, 12, 70,
			213, 228, 237, 60, 32, 55, 39, 222, 64, 108, 159, 53, 218, 145,
			55, 144, 121, 172, 1, 143, 20, 48, 12, 48, 105, 64, 184, 247,
			37, 111, 75, 93, 251, 96, 65, 165, 157, 93, 251, 163, 114, 127,
			33, 25, 86, 149, 46, 129, 167, 50, 172, 52, 106, 13, 181, 201,
			22, 22, 223, 80, 182, 57, 253, 240, 165, 231, 251, 76, 237, 46,
			30, 44, 181, 26, 5, 55, 108, 139, 194, 35, 190, 87, 24, 169,
			82, 168, 243, 14, 143, 120, 208, 230, 133, 251, 126, 184, 199, 252,
			103, 53, 165, 131, 40, 160, 66, 133, 49, 33, 139, 234, 66, 167,
			23, 186, 14, 26, 163, 35, 77, 94, 237, 115, 173, 18, 60, 199,
			58, 10, 65, 119, 146, 143, 231, 137, 65, 104, 234, 30, 79, 172,
			229, 46, 121, 163, 137, 4, 158, 62, 23, 50, 234, 168, 169, 99,
			22, 133, 109, 225, 12, 116, 100, 67, 91, 86, 11, 190, 183, 23,
			177, 232, 80, 21, 163, 78, 79, 246, 253, 27, 234, 43, 153, 187,
			168, 46, 34, 72, 234, 200, 137, 16, 49, 224, 109, 248, 120, 254,
			241, 210, 124, 127, 105, 222, 109, 206, 111, 221, 157, 223, 185, 59,
			223, 112, 230, 59, 79, 62, 118, 96, 219, 251, 138, 31, 120, 130,
			171, 226, 31, 1, 26, 173, 210, 80, 112, 205, 237, 65, 232, 50,
			229, 172, 31, 11, 120, 250, 188, 210, 168, 37, 169, 126, 83, 7,
			43, 55, 38, 23, 22, 159, 255, 100, 65, 95, 223, 197, 113, 238,
			203, 208, 213, 43, 129, 31, 75, 170, 138, 102, 3, 79, 45, 72,
			210, 170, 107, 107, 173, 107, 225, 36, 111, 101, 103, 34, 96, 126,
			181, 52, 191, 90, 34, 176, 136, 64, 134, 123, 234, 218, 140, 197,
			118, 74, 30, 65, 155, 13, 212, 6, 9, 59, 208, 229, 1, 143,
			152, 222, 106, 201, 54, 19, 58, 44, 167, 248, 59, 68, 253, 89,
			118, 198, 160, 214, 55, 83, 179, 228, 175, 13, 98, 219, 25, 51,
			67, 237, 159, 27, 230, 197, 236, 159, 26, 80, 31, 29, 251, 18,
			215, 15, 59, 202, 227, 21, 196, 194, 11, 218, 227, 165, 7, 121,
			115, 237, 1, 59, 67, 33, 209, 21, 190, 237, 172, 64, 222, 116,
			88, 120, 2, 94, 208, 246, 135, 194, 219, 199, 211, 211, 25, 50,
			129, 234, 77, 40, 253, 78, 37, 164, 129, 228, 212, 185, 132, 180,
			144, 164, 23, 200, 175, 180, 49, 6, 181, 255, 216, 48, 105, 246,
			63, 13, 168, 134, 193, 82, 192, 187, 250, 112, 120, 228, 136, 201,
			146, 163, 20, 158, 174, 222, 124, 196, 172, 198, 19, 211, 83, 215,
			62, 243, 135, 92, 232, 107, 186, 17, 51, 117, 153, 40, 164, 231,
			251, 208, 99, 251, 28, 130, 113, 153, 138, 117, 60, 145, 232, 35,
			141, 62, 181, 118, 194, 8, 79, 139, 201, 145, 250, 56, 96, 241,
			73, 42, 31, 255, 35, 111, 0, 197, 152, 80, 118, 38, 160, 24,
			202, 236, 169, 51, 9, 105, 33, 121, 126, 118, 111, 82, 135, 87,
			242, 175, 55, 201, 146, 23, 116, 34, 86, 96, 131, 1, 15, 186,
			94, 192, 11, 110, 20, 6, 124, 233, 235, 33, 231, 1, 122, 105,
			65, 240, 104, 223, 107, 199, 55, 240, 116, 70, 117, 63, 83, 221,
			217, 119, 189, 8, 228, 126, 110, 18, 90, 231, 131, 48, 146, 37,
			156, 86, 231, 95, 15, 185, 144, 244, 67, 66, 52, 155, 225, 208,
			115, 213, 107, 192, 116, 125, 90, 181, 180, 134, 158, 75, 31, 145,
			115, 126, 200, 220, 103, 113, 212, 14, 35, 253, 50, 48, 179, 234,
			56, 99, 210, 157, 147, 140, 157, 237, 144, 185, 149, 116, 86, 253,
			172, 127, 132, 166, 223, 39, 179, 154, 129, 203, 133, 10, 128, 94,
			24, 204, 89, 74, 252, 121, 213, 81, 26, 181, 83, 74, 236, 158,
			183, 207, 231, 108, 213, 175, 190, 179, 55, 201, 217, 163, 34, 232,
			117, 114, 218, 29, 202, 103, 184, 229, 218, 158, 60, 84, 198, 156,
			169, 207, 184, 67, 185, 17, 55, 229, 254, 197, 36, 23, 142, 232,
			42, 6, 97, 32, 56, 253, 156, 76, 10, 201, 228, 80, 191, 135,
			156, 93, 253, 248, 237, 214, 233, 25, 78, 67, 13, 175, 199, 211,
			142, 193, 104, 30, 135, 113, 131, 156, 227, 47, 6, 94, 164, 142,
			254, 207, 112, 105, 148, 173, 51, 171, 217, 227, 143, 42, 78, 154,
			130, 235, 103, 71, 83, 176, 145, 222, 32, 103, 152, 16, 94, 55,
			224, 238, 51, 119, 40, 197, 156, 13, 214, 194, 116, 253, 116, 210,
			88, 26, 74, 129, 131, 220, 136, 121, 129, 23, 116, 245, 160, 9,
			61, 40, 105, 196, 65, 185, 219, 100, 82, 235, 79, 103, 201, 153,
			86, 245, 71, 213, 218, 163, 234, 179, 114, 189, 94, 171, 159, 207,
			208, 73, 98, 214, 126, 116, 222, 160, 231, 201, 233, 164, 171, 213,
			170, 148, 206, 155, 185, 251, 232, 65, 62, 103, 130, 35, 151, 223,
			208, 131, 40, 177, 149, 30, 166, 210, 67, 125, 231, 222, 195, 85,
			24, 99, 164, 49, 205, 253, 149, 65, 104, 137, 183, 125, 22, 29,
			17, 240, 128, 156, 101, 251, 204, 243, 49, 144, 62, 75, 121, 205,
			172, 222, 56, 178, 72, 39, 39, 58, 165, 161, 172, 159, 73, 167,
			98, 79, 118, 137, 88, 165, 161, 68, 165, 2, 214, 231, 177, 182,
			234, 59, 117, 50, 115, 228, 100, 15, 236, 41, 227, 188, 57, 82,
			250, 136, 140, 88, 233, 11, 100, 118, 219, 19, 218, 59, 18, 201,
			185, 255, 54, 8, 29, 111, 141, 221, 236, 51, 50, 169, 84, 70,
			55, 67, 11, 230, 143, 88, 112, 114, 130, 163, 125, 46, 158, 148,
			253, 133, 65, 38, 84, 11, 61, 75, 204, 20, 107, 243, 205, 254,
			101, 126, 103, 255, 250, 46, 91, 50, 55, 75, 206, 41, 125, 71,
			112, 231, 254, 222, 32, 231, 71, 109, 177, 201, 183, 227, 229, 215,
			6, 95, 63, 105, 240, 216, 96, 181, 96, 106, 120, 246, 11, 189,
			78, 199, 237, 156, 39, 103, 71, 91, 0, 57, 197, 171, 149, 110,
			12, 13, 79, 150, 76, 37, 254, 174, 12, 152, 170, 167, 244, 234,
			63, 165, 32, 238, 146, 153, 177, 157, 77, 175, 189, 35, 162, 101,
			225, 93, 65, 65, 115, 76, 253, 250, 4, 199, 227, 91, 231, 4,
			199, 19, 91, 98, 149, 147, 217, 74, 176, 207, 3, 25, 70, 135,
			187, 250, 29, 42, 66, 49, 99, 158, 120, 76, 204, 201, 125, 112,
			76, 204, 27, 156, 120, 245, 111, 13, 114, 170, 18, 96, 93, 38,
			233, 14, 33, 35, 79, 164, 87, 223, 234, 162, 154, 247, 181, 119,
			184, 48, 189, 79, 166, 146, 117, 166, 87, 222, 178, 252, 154, 213,
			135, 223, 234, 28, 235, 215, 159, 92, 123, 71, 126, 124, 240, 235,
			143, 200, 36, 181, 237, 204, 55, 6, 249, 103, 67, 61, 19, 219,
			25, 186, 250, 75, 227, 200, 139, 239, 202, 154, 58, 136, 109, 183,
			54, 42, 80, 28, 202, 94, 24, 9, 231, 45, 207, 190, 45, 161,
			42, 183, 248, 113, 109, 244, 72, 234, 9, 232, 134, 251, 60, 10,
			240, 144, 26, 184, 241, 155, 95, 113, 192, 218, 200, 216, 107, 243,
			0, 203, 215, 135, 60, 18, 94, 24, 192, 170, 179, 156, 148, 22,
			186, 252, 238, 132, 195, 192, 77, 174, 182, 183, 43, 27, 229, 106,
			163, 12, 29, 207, 199, 218, 97, 154, 152, 86, 134, 90, 147, 153,
			197, 248, 105, 98, 42, 115, 49, 126, 28, 32, 153, 59, 201, 131,
			3, 126, 18, 98, 78, 102, 168, 125, 58, 115, 201, 192, 146, 113,
			18, 75, 198, 211, 83, 103, 200, 63, 24, 196, 158, 196, 146, 209,
			162, 102, 41, 251, 151, 170, 98, 76, 60, 21, 53, 111, 51, 223,
			215, 167, 47, 29, 86, 176, 148, 137, 212, 16, 240, 189, 125, 30,
			112, 161, 111, 252, 187, 92, 66, 169, 213, 36, 160, 247, 86, 31,
			75, 78, 60, 65, 53, 184, 126, 145, 173, 151, 139, 165, 157, 178,
			186, 218, 118, 185, 100, 158, 47, 240, 204, 37, 213, 189, 127, 32,
			177, 252, 26, 189, 77, 43, 73, 170, 18, 35, 241, 131, 172, 67,
			200, 105, 50, 49, 169, 138, 69, 139, 78, 206, 38, 148, 73, 45,
			74, 63, 74, 40, 139, 90, 180, 176, 78, 182, 149, 69, 6, 181,
			222, 51, 75, 217, 207, 97, 108, 163, 188, 221, 32, 53, 4, 194,
			131, 128, 71, 162, 231, 13, 112, 29, 75, 173, 166, 72, 229, 26,
			200, 46, 149, 139, 72, 191, 151, 202, 53, 44, 106, 189, 87, 88,
			87, 16, 27, 212, 158, 203, 92, 209, 16, 227, 156, 185, 169, 15,
			200, 30, 177, 39, 13, 68, 248, 178, 89, 202, 182, 96, 108, 71,
			129, 228, 190, 175, 15, 244, 113, 177, 6, 108, 47, 28, 74, 96,
			190, 175, 93, 137, 43, 53, 32, 77, 75, 170, 206, 214, 16, 163,
			226, 218, 132, 88, 75, 67, 161, 115, 57, 214, 210, 80, 232, 92,
			142, 181, 52, 20, 58, 151, 11, 235, 228, 47, 12, 98, 78, 154,
			212, 134, 204, 13, 35, 251, 11, 3, 226, 141, 156, 42, 16, 191,
			95, 11, 168, 239, 110, 136, 209, 83, 4, 214, 199, 251, 28, 60,
			61, 218, 11, 131, 130, 203, 247, 134, 221, 174, 23, 116, 29, 245,
			160, 32, 184, 158, 17, 87, 205, 233, 11, 10, 180, 195, 254, 128,
			73, 111, 207, 243, 61, 121, 8, 97, 132, 71, 207, 152, 232, 14,
			89, 196, 2, 201, 149, 9, 8, 25, 174, 26, 76, 157, 35, 51,
			196, 158, 52, 17, 178, 235, 102, 81, 233, 111, 42, 219, 174, 79,
			158, 79, 40, 147, 90, 215, 103, 115, 9, 101, 81, 235, 250, 210,
			231, 241, 52, 131, 90, 57, 243, 94, 220, 133, 139, 144, 155, 60,
			155, 80, 38, 181, 114, 231, 174, 38, 148, 69, 173, 220, 226, 26,
			46, 156, 157, 161, 246, 124, 230, 182, 145, 30, 167, 230, 167, 178,
			228, 79, 146, 227, 148, 181, 96, 206, 101, 127, 31, 70, 133, 11,
			58, 18, 46, 14, 150, 58, 144, 100, 19, 125, 58, 142, 221, 215,
			1, 168, 242, 131, 196, 199, 244, 13, 8, 1, 159, 35, 58, 42,
			66, 240, 254, 64, 30, 222, 3, 6, 1, 63, 208, 124, 14, 240,
			208, 177, 199, 223, 194, 79, 173, 177, 62, 61, 89, 11, 230, 84,
			66, 25, 212, 90, 152, 190, 144, 80, 22, 181, 22, 46, 189, 79,
			238, 197, 39, 39, 235, 19, 115, 62, 235, 192, 177, 154, 92, 221,
			51, 169, 31, 13, 116, 212, 235, 30, 115, 97, 143, 249, 44, 104,
			171, 181, 140, 89, 25, 147, 56, 251, 124, 66, 33, 175, 89, 72,
			40, 139, 90, 159, 220, 248, 136, 60, 84, 98, 76, 106, 229, 205,
			107, 217, 10, 156, 40, 7, 212, 53, 29, 244, 134, 125, 22, 64,
			39, 242, 120, 224, 250, 135, 48, 222, 31, 187, 120, 114, 31, 122,
			212, 80, 115, 2, 25, 39, 134, 162, 53, 249, 233, 108, 66, 89,
			212, 202, 127, 136, 235, 104, 219, 25, 43, 67, 237, 37, 115, 197,
			210, 125, 22, 66, 178, 68, 230, 136, 32, 147, 72, 225, 242, 45,
			219, 87, 178, 46, 140, 151, 251, 90, 53, 225, 169, 155, 90, 5,
			65, 138, 143, 142, 67, 163, 235, 182, 94, 120, 0, 125, 22, 28,
			18, 144, 161, 100, 190, 222, 144, 163, 48, 133, 81, 90, 12, 7,
			24, 17, 29, 66, 206, 146, 83, 90, 232, 4, 74, 29, 163, 13,
			106, 45, 207, 188, 63, 162, 45, 106, 45, 103, 47, 147, 63, 211,
			46, 102, 81, 235, 150, 73, 179, 127, 104, 0, 22, 146, 250, 132,
			169, 86, 103, 36, 135, 117, 121, 160, 238, 85, 61, 21, 198, 210,
			245, 43, 181, 154, 133, 120, 68, 167, 227, 5, 158, 60, 116, 136,
			214, 81, 157, 108, 5, 235, 243, 113, 166, 111, 118, 50, 79, 28,
			3, 223, 154, 64, 141, 18, 240, 45, 131, 90, 183, 166, 207, 36,
			20, 106, 123, 126, 86, 109, 27, 131, 218, 119, 50, 85, 189, 109,
			208, 73, 238, 76, 93, 38, 140, 216, 182, 138, 119, 107, 230, 197,
			108, 19, 244, 153, 39, 78, 26, 113, 176, 211, 77, 201, 242, 51,
			223, 119, 0, 42, 234, 126, 216, 235, 227, 48, 22, 168, 235, 180,
			118, 143, 183, 191, 34, 233, 149, 7, 240, 40, 194, 252, 171, 149,
			52, 204, 204, 36, 202, 152, 74, 40, 131, 90, 107, 211, 231, 18,
			202, 162, 214, 26, 189, 160, 60, 196, 192, 221, 125, 215, 92, 215,
			30, 98, 168, 253, 125, 247, 212, 25, 242, 7, 38, 153, 68, 18,
			117, 253, 220, 190, 148, 253, 181, 1, 71, 142, 55, 201, 197, 101,
			16, 202, 244, 119, 54, 65, 24, 245, 153, 239, 31, 166, 10, 171,
			21, 226, 29, 54, 244, 37, 137, 49, 246, 58, 227, 86, 122, 2,
			212, 239, 103, 130, 46, 6, 191, 97, 240, 85, 16, 30, 4, 14,
			28, 189, 191, 212, 83, 72, 26, 133, 135, 130, 139, 56, 54, 240,
			96, 216, 143, 25, 167, 25, 178, 237, 123, 106, 195, 132, 92, 40,
			237, 144, 39, 137, 115, 199, 33, 143, 111, 250, 227, 65, 106, 197,
			135, 130, 143, 107, 170, 249, 197, 254, 106, 196, 113, 228, 115, 123,
			118, 68, 155, 212, 250, 252, 226, 123, 228, 76, 140, 144, 65, 173,
			31, 218, 51, 105, 183, 161, 232, 201, 17, 109, 82, 235, 135, 211,
			36, 29, 110, 82, 171, 104, 191, 151, 118, 227, 244, 162, 125, 126,
			68, 99, 255, 133, 139, 228, 111, 12, 229, 42, 6, 181, 54, 205,
			185, 236, 159, 27, 223, 53, 194, 86, 58, 227, 51, 14, 152, 64,
			0, 101, 82, 43, 69, 186, 82, 140, 127, 244, 213, 241, 184, 239,
			106, 48, 226, 95, 94, 245, 116, 44, 230, 122, 143, 104, 132, 195,
			136, 224, 82, 135, 250, 119, 115, 169, 167, 25, 19, 168, 98, 226,
			105, 104, 253, 102, 28, 116, 13, 21, 13, 55, 47, 189, 79, 54,
			149, 45, 38, 181, 182, 204, 229, 236, 26, 28, 59, 97, 161, 61,
			234, 6, 125, 60, 224, 141, 106, 37, 61, 156, 143, 124, 219, 156,
			68, 70, 151, 19, 202, 160, 214, 214, 149, 239, 39, 148, 69, 173,
			45, 167, 64, 126, 168, 36, 90, 212, 122, 96, 126, 148, 189, 9,
			71, 142, 251, 42, 200, 143, 234, 135, 111, 73, 41, 134, 105, 217,
			200, 34, 165, 38, 168, 245, 96, 102, 54, 161, 12, 106, 61, 160,
			215, 18, 10, 133, 229, 110, 144, 72, 73, 182, 169, 181, 99, 126,
			148, 69, 118, 99, 119, 8, 71, 37, 31, 43, 234, 226, 29, 165,
			38, 56, 144, 68, 51, 146, 60, 65, 48, 16, 195, 61, 92, 194,
			176, 115, 212, 156, 84, 87, 91, 9, 77, 169, 9, 106, 237, 164,
			186, 218, 6, 181, 118, 82, 93, 109, 139, 90, 59, 185, 27, 42,
			76, 153, 212, 222, 205, 60, 212, 97, 10, 177, 220, 157, 202, 146,
			31, 16, 219, 86, 53, 70, 221, 156, 203, 22, 190, 155, 235, 105,
			249, 166, 10, 243, 245, 216, 47, 116, 137, 82, 143, 253, 66, 23,
			37, 245, 75, 239, 147, 167, 74, 142, 65, 173, 150, 121, 57, 91,
			5, 5, 209, 168, 230, 76, 227, 8, 110, 99, 22, 232, 16, 135,
			225, 128, 33, 126, 105, 199, 72, 11, 242, 6, 53, 12, 27, 185,
			167, 212, 4, 181, 90, 49, 40, 186, 0, 106, 209, 75, 9, 101,
			81, 171, 245, 65, 22, 79, 6, 136, 207, 163, 204, 85, 133, 9,
			174, 242, 163, 169, 203, 10, 43, 155, 218, 143, 51, 158, 198, 10,
			17, 125, 60, 149, 37, 191, 139, 223, 211, 212, 122, 106, 158, 201,
			174, 106, 19, 48, 101, 240, 65, 196, 213, 235, 140, 3, 234, 244,
			115, 244, 226, 5, 139, 69, 201, 25, 238, 162, 25, 98, 219, 246,
			116, 134, 90, 79, 103, 78, 43, 77, 236, 105, 4, 107, 140, 50,
			53, 165, 132, 18, 106, 253, 216, 164, 122, 18, 201, 80, 235, 199,
			202, 24, 219, 182, 49, 211, 255, 212, 116, 117, 28, 183, 85, 166,
			255, 41, 57, 67, 174, 147, 73, 164, 112, 45, 159, 219, 23, 179,
			52, 253, 61, 101, 236, 133, 113, 156, 179, 227, 188, 252, 220, 30,
			163, 13, 106, 61, 159, 57, 55, 162, 45, 106, 61, 167, 23, 200,
			31, 25, 49, 79, 131, 90, 109, 251, 98, 86, 142, 231, 208, 49,
			206, 240, 27, 38, 228, 166, 30, 175, 202, 142, 49, 143, 98, 241,
			182, 120, 83, 170, 30, 211, 26, 87, 180, 61, 166, 53, 174, 105,
			123, 76, 107, 92, 213, 54, 189, 128, 199, 88, 219, 182, 17, 135,
			158, 153, 203, 254, 163, 113, 98, 65, 112, 135, 37, 191, 126, 85,
			219, 179, 207, 92, 126, 228, 116, 145, 28, 41, 148, 91, 226, 209,
			140, 121, 129, 24, 63, 213, 129, 23, 232, 215, 7, 44, 224, 208,
			94, 22, 3, 161, 248, 197, 129, 83, 223, 123, 143, 126, 108, 27,
			87, 29, 68, 111, 124, 238, 170, 35, 163, 203, 125, 62, 10, 178,
			182, 153, 177, 81, 239, 148, 154, 164, 86, 111, 230, 108, 66, 25,
			212, 234, 157, 251, 48, 161, 44, 106, 245, 64, 253, 248, 13, 35,
			192, 151, 177, 23, 79, 24, 212, 250, 114, 234, 178, 106, 158, 164,
			150, 159, 185, 162, 154, 39, 13, 106, 249, 83, 31, 40, 231, 62,
			69, 237, 126, 70, 106, 231, 62, 101, 80, 171, 63, 149, 85, 174,
			117, 10, 93, 43, 48, 35, 237, 90, 167, 148, 107, 5, 228, 156,
			74, 104, 167, 180, 107, 133, 54, 85, 128, 159, 138, 221, 40, 140,
			23, 228, 84, 236, 70, 225, 204, 153, 17, 109, 81, 43, 60, 63,
			155, 78, 55, 168, 53, 176, 87, 211, 110, 44, 174, 7, 246, 135,
			35, 26, 251, 175, 46, 141, 104, 139, 90, 131, 229, 149, 116, 186,
			73, 173, 175, 237, 235, 105, 55, 86, 198, 95, 143, 73, 71, 246,
			95, 207, 92, 25, 209, 22, 181, 190, 190, 6, 106, 3, 157, 66,
			213, 133, 121, 69, 219, 165, 48, 22, 49, 198, 167, 20, 198, 98,
			230, 124, 66, 25, 212, 18, 179, 239, 39, 148, 69, 45, 145, 213,
			96, 78, 81, 107, 63, 147, 85, 160, 77, 25, 212, 218, 159, 122,
			95, 129, 57, 77, 237, 131, 204, 55, 26, 204, 105, 131, 90, 7,
			83, 115, 10, 204, 105, 4, 243, 133, 249, 51, 13, 230, 180, 2,
			243, 5, 57, 163, 204, 153, 214, 96, 30, 198, 96, 78, 199, 96,
			30, 198, 230, 76, 199, 96, 30, 198, 96, 78, 199, 96, 30, 198,
			96, 78, 107, 48, 95, 218, 87, 211, 110, 220, 28, 47, 199, 166,
			35, 152, 47, 103, 62, 24, 209, 22, 181, 94, 94, 249, 48, 157,
			110, 82, 235, 149, 125, 41, 237, 70, 48, 95, 217, 83, 35, 218,
			160, 214, 171, 233, 217, 17, 109, 81, 235, 213, 197, 247, 20, 152,
			211, 168, 250, 107, 115, 78, 219, 165, 192, 124, 29, 131, 57, 173,
			192, 124, 29, 59, 236, 180, 50, 226, 245, 185, 11, 9, 101, 81,
			235, 245, 165, 247, 147, 231, 155, 255, 11, 0, 0, 255, 255, 24,
			52, 37, 199, 149, 49, 0, 0},
	)
}

// FileDescriptorSet returns a descriptor set for this proto package, which
// includes all defined services, and all transitive dependencies.
//
// Will not return nil.
//
// Do NOT modify the returned descriptor.
func FileDescriptorSet() *descriptorpb.FileDescriptorSet {
	// We just need ONE of the service names to look up the FileDescriptorSet.
	ret, err := discovery.GetDescriptorSet("drone_queen.Drone")
	if err != nil {
		panic(err)
	}
	return ret
}
