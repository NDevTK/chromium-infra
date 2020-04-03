// Code generated by cproto. DO NOT EDIT.

package api

import discovery "go.chromium.org/luci/grpc/discovery"

import "github.com/golang/protobuf/protoc-gen-go/descriptor"

func init() {
	discovery.RegisterDescriptorSetCompressed(
		[]string{
			"fleet.Configuration", "fleet.Registration",
		},
		[]byte{31, 139,
			8, 0, 0, 0, 0, 0, 0, 255, 228, 91, 75, 115, 27, 59,
			118, 102, 3, 32, 213, 4, 223, 32, 41, 146, 173, 23, 68, 189,
			45, 89, 242, 149, 237, 177, 37, 123, 238, 181, 36, 203, 51, 246,
			189, 126, 12, 229, 113, 85, 102, 163, 180, 72, 72, 234, 186, 20,
			91, 211, 221, 210, 29, 103, 53, 201, 46, 217, 100, 157, 44, 83,
			149, 217, 37, 153, 69, 22, 249, 5, 169, 212, 44, 242, 31, 50,
			73, 37, 63, 33, 85, 89, 164, 82, 7, 13, 52, 65, 73, 126,
			100, 82, 89, 76, 221, 149, 136, 6, 240, 225, 224, 224, 224, 224,
			124, 7, 16, 253, 79, 155, 86, 187, 167, 129, 127, 38, 14, 187,
			254, 224, 216, 59, 89, 63, 15, 252, 200, 103, 233, 227, 190, 16,
			81, 251, 47, 45, 90, 220, 147, 213, 111, 250, 110, 116, 236, 7,
			103, 108, 137, 34, 175, 215, 180, 184, 181, 156, 219, 108, 172, 203,
			102, 235, 163, 77, 158, 63, 237, 32, 175, 199, 90, 212, 190, 60,
			59, 12, 251, 126, 20, 54, 17, 183, 150, 211, 157, 177, 203, 179,
			3, 40, 178, 54, 205, 159, 185, 131, 139, 99, 183, 27, 93, 4,
			34, 104, 98, 110, 45, 103, 59, 35, 223, 24, 167, 185, 158, 8,
			187, 129, 119, 30, 121, 254, 160, 73, 100, 19, 243, 83, 123, 143,
			102, 95, 31, 188, 19, 65, 232, 249, 3, 86, 163, 233, 75, 183,
			127, 33, 164, 100, 217, 78, 92, 184, 10, 130, 174, 131, 44, 211,
			242, 85, 233, 111, 198, 218, 157, 252, 153, 227, 13, 142, 3, 119,
			163, 239, 29, 133, 27, 114, 226, 27, 82, 89, 225, 198, 137, 255,
			226, 87, 25, 154, 97, 132, 164, 90, 22, 253, 27, 139, 90, 121,
			134, 73, 138, 109, 254, 149, 197, 247, 252, 243, 247, 129, 119, 114,
			26, 241, 205, 59, 155, 119, 248, 219, 83, 193, 229, 128, 222, 197,
			25, 127, 125, 192, 119, 46, 162, 83, 63, 8, 215, 249, 78, 191,
			207, 101, 187, 144, 7, 34, 20, 193, 165, 232, 173, 83, 254, 211,
			80, 112, 255, 152, 71, 167, 94, 200, 67, 255, 34, 232, 10, 222,
			245, 123, 130, 123, 33, 63, 241, 47, 69, 48, 16, 61, 126, 244,
			158, 187, 124, 247, 224, 233, 237, 48, 122, 223, 23, 188, 239, 117,
			197, 32, 20, 60, 58, 117, 35, 222, 117, 7, 252, 72, 80, 126,
			236, 95, 12, 122, 220, 27, 240, 232, 84, 240, 111, 158, 239, 237,
			191, 58, 216, 231, 199, 94, 95, 172, 83, 106, 83, 11, 49, 156,
			177, 11, 240, 203, 102, 56, 155, 250, 130, 46, 82, 100, 231, 228,
			79, 103, 130, 119, 196, 207, 47, 188, 32, 30, 73, 206, 184, 123,
			251, 68, 12, 110, 159, 248, 148, 82, 138, 72, 138, 145, 92, 170,
			108, 81, 74, 49, 73, 89, 12, 231, 236, 113, 154, 163, 132, 164,
			80, 138, 225, 60, 153, 164, 5, 154, 134, 2, 97, 36, 79, 114,
			77, 154, 143, 139, 25, 168, 172, 233, 146, 197, 112, 190, 222, 208,
			37, 204, 112, 222, 153, 80, 40, 22, 195, 5, 210, 80, 40, 22,
			97, 164, 64, 242, 147, 170, 165, 149, 134, 202, 172, 46, 65, 83,
			202, 116, 9, 51, 92, 168, 107, 89, 16, 195, 197, 68, 22, 68,
			24, 41, 146, 130, 30, 15, 165, 161, 146, 234, 146, 197, 112, 49,
			151, 212, 97, 134, 139, 137, 44, 152, 225, 18, 153, 80, 40, 152,
			48, 82, 34, 69, 45, 11, 78, 67, 165, 70, 193, 22, 195, 165,
			220, 184, 46, 65, 199, 150, 35, 245, 101, 49, 194, 82, 245, 88,
			95, 32, 49, 179, 43, 18, 221, 2, 125, 85, 73, 93, 162, 91,
			82, 95, 85, 194, 170, 18, 193, 66, 169, 52, 84, 82, 93, 178,
			24, 174, 230, 202, 186, 132, 25, 174, 86, 107, 10, 197, 98, 184,
			166, 100, 180, 164, 190, 106, 164, 90, 87, 45, 65, 95, 181, 4,
			5, 70, 175, 41, 25, 45, 169, 175, 154, 146, 17, 49, 210, 0,
			91, 6, 25, 1, 176, 97, 55, 37, 58, 2, 25, 155, 74, 70,
			36, 101, 108, 146, 134, 35, 17, 144, 148, 177, 169, 208, 145, 148,
			177, 169, 100, 68, 82, 198, 102, 181, 118, 148, 145, 22, 116, 151,
			254, 35, 165, 15, 227, 253, 228, 158, 159, 139, 193, 137, 55, 16,
			27, 23, 3, 239, 216, 19, 189, 219, 241, 230, 114, 207, 189, 141,
			203, 47, 54, 98, 191, 116, 17, 184, 176, 95, 71, 220, 147, 115,
			147, 235, 106, 239, 211, 201, 231, 103, 231, 126, 16, 141, 238, 235,
			16, 172, 88, 132, 17, 91, 160, 197, 190, 223, 117, 251, 135, 96,
			254, 231, 110, 116, 170, 182, 121, 65, 126, 125, 166, 62, 182, 255,
			212, 162, 83, 31, 192, 9, 207, 253, 65, 40, 216, 93, 154, 57,
			119, 195, 80, 128, 55, 196, 203, 185, 205, 137, 27, 189, 97, 71,
			132, 23, 253, 168, 163, 154, 66, 167, 99, 215, 235, 139, 94, 19,
			125, 70, 167, 184, 105, 251, 152, 214, 110, 170, 103, 95, 80, 251,
			92, 125, 81, 30, 185, 126, 51, 92, 210, 140, 77, 208, 172, 8,
			2, 63, 56, 60, 11, 79, 148, 63, 180, 229, 135, 151, 225, 201,
			102, 72, 11, 123, 166, 182, 217, 17, 173, 223, 168, 3, 54, 167,
			198, 249, 152, 166, 157, 249, 143, 55, 138, 213, 184, 155, 254, 25,
			118, 207, 189, 23, 127, 171, 28, 232, 244, 239, 191, 3, 181, 83,
			14, 205, 198, 14, 84, 253, 196, 41, 134, 233, 216, 36, 108, 173,
			76, 138, 145, 124, 170, 36, 183, 86, 70, 186, 61, 187, 78, 191,
			164, 36, 35, 221, 101, 17, 29, 57, 95, 240, 27, 21, 198, 61,
			249, 53, 228, 177, 209, 115, 189, 164, 225, 186, 220, 112, 153, 216,
			137, 22, 51, 19, 186, 4, 30, 111, 242, 161, 46, 129, 27, 219,
			251, 67, 237, 176, 43, 169, 218, 208, 97, 87, 236, 121, 250, 133,
			118, 216, 85, 52, 233, 204, 243, 131, 139, 115, 24, 76, 141, 233,
			13, 78, 248, 113, 224, 159, 113, 185, 69, 244, 180, 135, 110, 189,
			138, 42, 139, 218, 117, 131, 155, 66, 182, 225, 214, 171, 89, 211,
			173, 87, 157, 9, 237, 4, 199, 181, 131, 1, 55, 52, 110, 47,
			12, 157, 96, 3, 173, 106, 207, 70, 160, 148, 120, 189, 12, 195,
			141, 220, 140, 225, 3, 27, 124, 209, 240, 129, 141, 149, 91, 67,
			31, 216, 76, 64, 44, 2, 165, 196, 233, 101, 192, 45, 205, 24,
			46, 176, 153, 128, 128, 11, 108, 174, 220, 210, 46, 112, 2, 172,
			81, 187, 192, 9, 123, 114, 232, 2, 39, 209, 188, 225, 2, 39,
			209, 196, 180, 118, 115, 25, 168, 28, 55, 92, 224, 100, 99, 198,
			112, 129, 147, 237, 57, 133, 98, 49, 60, 133, 26, 10, 5, 220,
			244, 20, 154, 156, 87, 45, 193, 77, 79, 41, 45, 34, 41, 227,
			84, 150, 233, 18, 102, 120, 170, 62, 158, 56, 210, 127, 33, 180,
			30, 155, 132, 31, 222, 20, 196, 221, 58, 165, 44, 54, 165, 215,
			7, 79, 197, 165, 215, 21, 111, 223, 159, 11, 198, 104, 241, 233,
			254, 187, 231, 123, 251, 135, 207, 95, 189, 219, 249, 230, 249, 211,
			114, 138, 213, 105, 69, 125, 219, 251, 113, 231, 245, 203, 253, 221,
			215, 175, 191, 46, 91, 198, 231, 111, 118, 118, 15, 222, 238, 188,
			125, 254, 250, 85, 25, 177, 50, 205, 171, 207, 7, 251, 157, 119,
			175, 203, 248, 19, 33, 210, 175, 113, 188, 195, 203, 191, 255, 59,
			252, 127, 25, 34, 165, 141, 16, 41, 45, 67, 164, 116, 3, 172,
			32, 173, 66, 36, 185, 65, 210, 42, 10, 34, 76, 151, 16, 196,
			68, 227, 170, 161, 140, 130, 28, 85, 37, 3, 29, 82, 215, 37,
			196, 112, 161, 217, 82, 13, 101, 160, 163, 27, 202, 88, 38, 105,
			40, 235, 146, 134, 50, 150, 209, 85, 50, 92, 33, 101, 93, 66,
			12, 151, 140, 195, 250, 215, 99, 241, 161, 121, 237, 8, 110, 255,
			153, 69, 237, 111, 84, 13, 43, 83, 220, 119, 143, 212, 121, 10,
			63, 33, 148, 118, 189, 176, 47, 212, 81, 19, 23, 160, 93, 224,
			127, 167, 194, 126, 248, 201, 24, 37, 129, 219, 253, 86, 133, 249,
			242, 55, 244, 13, 79, 69, 255, 184, 153, 142, 251, 202, 2, 115,
			168, 125, 238, 135, 158, 140, 231, 51, 241, 249, 165, 203, 159, 176,
			191, 127, 77, 199, 246, 87, 251, 126, 217, 95, 43, 246, 248, 249,
			84, 205, 114, 10, 252, 149, 248, 69, 196, 223, 186, 39, 219, 252,
			1, 77, 14, 128, 188, 93, 166, 27, 250, 0, 40, 162, 138, 211,
			230, 122, 77, 227, 169, 9, 14, 209, 75, 164, 165, 235, 187, 71,
			166, 251, 47, 162, 60, 51, 220, 127, 113, 196, 253, 23, 179, 121,
			195, 253, 23, 75, 229, 97, 84, 95, 66, 85, 35, 170, 47, 161,
			98, 197, 136, 234, 75, 9, 10, 24, 123, 41, 91, 52, 162, 250,
			82, 133, 13, 163, 250, 50, 170, 24, 81, 125, 25, 149, 170, 70,
			84, 95, 78, 80, 96, 192, 114, 34, 11, 24, 127, 57, 145, 5,
			51, 92, 65, 204, 136, 234, 43, 168, 92, 49, 162, 250, 74, 130,
			2, 219, 164, 146, 45, 24, 81, 125, 165, 92, 81, 40, 132, 97,
			150, 204, 136, 16, 70, 24, 170, 104, 189, 144, 52, 84, 106, 20,
			2, 81, 127, 50, 35, 130, 25, 102, 201, 140, 228, 249, 57, 174,
			80, 210, 242, 112, 101, 122, 70, 233, 145, 195, 53, 45, 15, 87,
			45, 103, 26, 14, 215, 90, 61, 217, 178, 127, 71, 104, 249, 204,
			237, 158, 122, 3, 113, 232, 245, 70, 55, 237, 44, 205, 190, 140,
			171, 126, 71, 182, 251, 155, 239, 173, 43, 191, 129, 237, 86, 232,
			159, 88, 122, 243, 148, 72, 221, 185, 224, 59, 252, 98, 224, 253,
			252, 66, 240, 231, 79, 249, 177, 31, 112, 181, 14, 219, 148, 115,
			126, 139, 239, 192, 86, 122, 235, 158, 200, 42, 125, 50, 235, 54,
			124, 185, 39, 207, 232, 149, 184, 237, 129, 8, 60, 183, 207, 7,
			23, 103, 71, 34, 48, 58, 12, 155, 75, 229, 5, 43, 230, 126,
			44, 145, 92, 213, 216, 143, 67, 78, 154, 146, 156, 180, 108, 236,
			71, 211, 201, 255, 247, 20, 173, 156, 139, 192, 59, 63, 21, 129,
			219, 15, 63, 131, 106, 253, 189, 69, 241, 215, 239, 94, 178, 73,
			35, 39, 148, 87, 65, 255, 215, 239, 94, 170, 68, 16, 163, 100,
			224, 158, 233, 35, 64, 254, 102, 51, 52, 119, 230, 118, 15, 221,
			94, 47, 16, 97, 168, 78, 2, 122, 230, 118, 119, 226, 47, 35,
			212, 134, 124, 30, 181, 89, 162, 37, 247, 210, 245, 250, 238, 81,
			95, 28, 202, 80, 89, 158, 28, 233, 78, 49, 249, 252, 6, 190,
			182, 159, 208, 60, 72, 55, 136, 68, 112, 236, 118, 197, 167, 197,
			7, 48, 45, 62, 252, 110, 255, 177, 69, 113, 231, 205, 205, 19,
			239, 188, 249, 63, 77, 252, 134, 89, 144, 15, 205, 2, 134, 250,
			232, 44, 70, 100, 185, 54, 139, 99, 154, 57, 248, 206, 139, 186,
			167, 108, 198, 232, 91, 82, 125, 227, 170, 143, 76, 229, 6, 73,
			241, 141, 146, 62, 163, 37, 5, 150, 8, 251, 57, 3, 94, 147,
			247, 175, 45, 74, 158, 6, 110, 151, 77, 25, 189, 11, 170, 55,
			84, 168, 190, 87, 116, 140, 174, 233, 120, 135, 150, 67, 57, 216,
			161, 167, 37, 146, 162, 231, 54, 199, 71, 101, 209, 181, 157, 82,
			120, 101, 2, 16, 134, 184, 97, 248, 157, 31, 244, 84, 208, 146,
			148, 219, 83, 52, 45, 205, 231, 102, 215, 10, 213, 114, 93, 62,
			80, 205, 169, 173, 53, 241, 129, 22, 211, 52, 19, 207, 246, 119,
			242, 221, 127, 62, 30, 251, 238, 131, 239, 151, 239, 86, 44, 60,
			23, 179, 112, 112, 227, 5, 157, 132, 3, 183, 88, 176, 115, 116,
			101, 232, 197, 43, 206, 36, 15, 165, 239, 61, 84, 190, 215, 15,
			84, 12, 20, 185, 39, 163, 206, 182, 80, 48, 82, 154, 165, 36,
			25, 41, 157, 45, 53, 131, 159, 82, 169, 172, 198, 128, 40, 130,
			212, 156, 73, 254, 84, 167, 162, 47, 5, 135, 77, 197, 189, 99,
			62, 16, 162, 39, 122, 212, 8, 141, 42, 164, 100, 134, 70, 149,
			196, 161, 67, 104, 84, 201, 149, 140, 208, 168, 194, 170, 195, 208,
			136, 37, 169, 74, 8, 141, 24, 169, 212, 140, 208, 136, 141, 36,
			60, 89, 146, 170, 132, 104, 136, 181, 156, 97, 104, 84, 37, 220,
			8, 141, 170, 132, 77, 232, 240, 39, 3, 149, 204, 8, 141, 170,
			213, 9, 35, 52, 170, 78, 207, 12, 67, 163, 26, 153, 54, 66,
			163, 26, 169, 114, 35, 52, 170, 37, 90, 131, 208, 168, 70, 91,
			70, 104, 84, 155, 156, 162, 235, 113, 198, 160, 153, 154, 176, 156,
			54, 239, 136, 99, 17, 240, 200, 231, 254, 64, 112, 153, 166, 240,
			143, 185, 203, 79, 188, 75, 49, 224, 95, 191, 123, 73, 147, 172,
			66, 211, 174, 13, 179, 10, 45, 82, 49, 82, 171, 45, 210, 28,
			55, 210, 10, 45, 37, 65, 156, 86, 104, 169, 117, 139, 211, 10,
			45, 21, 40, 202, 180, 130, 67, 106, 70, 106, 213, 33, 173, 138,
			145, 90, 117, 70, 82, 171, 142, 90, 153, 56, 175, 224, 176, 170,
			206, 43, 76, 165, 22, 135, 121, 133, 41, 101, 121, 50, 175, 48,
			243, 89, 150, 23, 103, 29, 102, 200, 84, 193, 200, 58, 204, 168,
			25, 196, 89, 135, 25, 53, 131, 56, 235, 48, 163, 44, 79, 46,
			245, 236, 103, 89, 94, 156, 147, 152, 37, 51, 21, 35, 39, 49,
			155, 36, 119, 97, 126, 179, 106, 126, 113, 78, 98, 86, 89, 158,
			100, 153, 109, 101, 121, 72, 90, 94, 155, 204, 214, 84, 75, 176,
			188, 118, 130, 2, 226, 180, 115, 58, 91, 2, 198, 214, 110, 57,
			244, 189, 68, 193, 12, 207, 145, 233, 118, 159, 191, 184, 8, 35,
			25, 250, 4, 162, 235, 7, 61, 126, 42, 2, 177, 166, 252, 5,
			119, 123, 61, 209, 227, 125, 55, 18, 129, 49, 3, 254, 214, 135,
			202, 184, 131, 232, 109, 243, 51, 191, 39, 250, 107, 220, 188, 97,
			90, 227, 238, 217, 185, 8, 220, 19, 177, 198, 47, 253, 126, 228,
			158, 8, 61, 115, 176, 243, 57, 210, 158, 80, 114, 1, 5, 152,
			75, 180, 11, 118, 62, 167, 44, 20, 73, 59, 159, 83, 22, 138,
			25, 89, 73, 173, 125, 210, 66, 59, 111, 148, 133, 2, 210, 138,
			178, 80, 12, 171, 127, 75, 89, 40, 150, 235, 123, 139, 172, 196,
			154, 193, 114, 125, 111, 41, 9, 176, 92, 223, 91, 106, 125, 177,
			92, 223, 91, 202, 66, 49, 40, 116, 85, 89, 40, 150, 43, 184,
			74, 110, 85, 84, 75, 88, 193, 85, 165, 123, 44, 87, 112, 85,
			173, 32, 150, 43, 184, 170, 44, 148, 48, 178, 158, 250, 65, 108,
			161, 176, 31, 215, 237, 162, 180, 30, 2, 50, 222, 145, 214, 243,
			41, 11, 37, 114, 6, 119, 200, 122, 28, 108, 18, 57, 131, 59,
			36, 175, 75, 22, 195, 119, 10, 37, 93, 194, 12, 223, 97, 85,
			53, 134, 197, 240, 230, 103, 89, 40, 145, 243, 219, 36, 119, 106,
			10, 7, 230, 183, 169, 230, 71, 228, 252, 54, 115, 122, 12, 152,
			223, 38, 171, 210, 99, 57, 6, 98, 248, 46, 153, 110, 255, 193,
			255, 151, 109, 105, 1, 193, 248, 239, 146, 77, 45, 32, 24, 255,
			93, 181, 140, 241, 94, 188, 171, 12, 41, 182, 247, 187, 147, 83,
			244, 11, 10, 30, 145, 108, 165, 30, 91, 206, 194, 199, 13, 41,
			142, 73, 98, 91, 2, 50, 184, 101, 203, 172, 18, 73, 195, 58,
			109, 43, 43, 72, 203, 149, 216, 38, 91, 241, 64, 105, 185, 18,
			219, 106, 37, 210, 114, 37, 182, 213, 74, 164, 229, 74, 108, 171,
			125, 156, 6, 249, 30, 37, 40, 160, 235, 71, 100, 187, 166, 90,
			130, 174, 31, 41, 93, 167, 165, 174, 31, 229, 52, 10, 232, 250,
			145, 178, 165, 12, 35, 95, 166, 246, 98, 91, 202, 88, 12, 127,
			105, 231, 37, 122, 6, 100, 252, 138, 196, 164, 59, 35, 101, 252,
			138, 124, 25, 83, 226, 140, 148, 241, 43, 133, 158, 145, 50, 126,
			149, 43, 232, 18, 102, 248, 43, 69, 186, 51, 32, 227, 19, 229,
			107, 50, 82, 198, 39, 228, 43, 166, 90, 130, 140, 79, 18, 20,
			144, 241, 137, 242, 53, 25, 41, 227, 19, 117, 202, 101, 192, 30,
			118, 200, 138, 66, 129, 69, 219, 33, 79, 38, 84, 75, 148, 129,
			202, 170, 46, 89, 12, 239, 212, 230, 117, 9, 51, 188, 179, 180,
			172, 80, 48, 195, 187, 164, 169, 80, 192, 135, 236, 146, 157, 21,
			213, 18, 124, 200, 110, 34, 11, 236, 252, 221, 156, 198, 4, 31,
			178, 59, 222, 144, 250, 26, 99, 100, 63, 245, 163, 88, 95, 99,
			22, 195, 251, 118, 65, 162, 143, 129, 190, 158, 169, 139, 183, 49,
			169, 175, 103, 100, 63, 214, 248, 152, 164, 121, 207, 20, 250, 152,
			212, 215, 51, 69, 243, 198, 164, 190, 158, 85, 107, 18, 221, 102,
			228, 121, 234, 235, 24, 221, 182, 24, 126, 174, 208, 109, 64, 127,
			161, 208, 109, 137, 254, 130, 60, 143, 209, 109, 137, 254, 66, 161,
			219, 18, 253, 133, 66, 183, 37, 250, 11, 133, 158, 101, 228, 101,
			234, 117, 140, 158, 181, 24, 126, 105, 199, 94, 41, 11, 232, 175,
			20, 122, 86, 162, 191, 34, 47, 227, 85, 202, 74, 244, 87, 10,
			61, 43, 209, 95, 41, 244, 172, 68, 127, 165, 208, 41, 35, 63,
			129, 160, 21, 208, 169, 197, 240, 79, 236, 162, 68, 167, 128, 222,
			81, 232, 84, 162, 119, 200, 79, 98, 4, 42, 209, 59, 10, 157,
			74, 244, 78, 46, 169, 195, 12, 119, 12, 2, 252, 15, 14, 45,
			40, 94, 253, 105, 242, 235, 220, 156, 116, 119, 174, 228, 73, 157,
			107, 73, 24, 231, 58, 201, 110, 255, 187, 69, 199, 84, 70, 134,
			113, 131, 220, 148, 21, 29, 73, 178, 53, 146, 223, 172, 82, 91,
			15, 34, 201, 205, 144, 66, 233, 172, 93, 39, 105, 192, 126, 72,
			139, 74, 122, 37, 136, 98, 58, 181, 17, 58, 173, 6, 248, 113,
			170, 83, 232, 154, 31, 216, 30, 45, 39, 243, 212, 0, 100, 132,
			42, 233, 12, 198, 16, 162, 164, 123, 168, 79, 187, 54, 205, 196,
			89, 141, 246, 63, 97, 90, 24, 25, 49, 97, 150, 150, 193, 44,
			77, 242, 143, 62, 143, 252, 127, 146, 87, 183, 168, 61, 240, 186,
			135, 114, 172, 152, 176, 141, 13, 188, 238, 43, 24, 238, 33, 45,
			124, 123, 121, 102, 112, 193, 180, 28, 179, 106, 164, 2, 18, 34,
			152, 255, 246, 242, 108, 200, 2, 31, 210, 66, 112, 110, 246, 204,
			140, 244, 52, 249, 121, 39, 31, 156, 27, 61, 223, 208, 230, 64,
			68, 223, 249, 193, 183, 135, 177, 106, 12, 144, 177, 143, 82, 209,
			113, 213, 47, 190, 203, 49, 41, 53, 233, 5, 110, 183, 105, 203,
			222, 57, 131, 22, 119, 100, 5, 91, 165, 149, 158, 56, 239, 251,
			239, 207, 196, 32, 58, 140, 188, 238, 183, 34, 106, 102, 165, 42,
			202, 195, 138, 183, 242, 251, 213, 151, 51, 244, 250, 203, 153, 255,
			178, 104, 233, 202, 218, 179, 37, 90, 10, 224, 196, 18, 131, 174,
			56, 60, 242, 221, 160, 167, 214, 181, 152, 124, 222, 133, 175, 108,
			150, 230, 143, 46, 188, 126, 239, 48, 114, 131, 19, 161, 105, 126,
			78, 126, 123, 43, 63, 1, 183, 149, 39, 171, 90, 203, 184, 192,
			182, 65, 46, 169, 175, 232, 253, 121, 188, 146, 197, 205, 214, 21,
			83, 28, 94, 115, 117, 104, 111, 120, 229, 181, 76, 203, 39, 190,
			127, 210, 135, 141, 220, 19, 177, 41, 196, 119, 11, 197, 248, 251,
			158, 223, 19, 175, 110, 200, 210, 100, 174, 90, 211, 39, 40, 246,
			127, 84, 98, 138, 253, 248, 123, 75, 177, 225, 167, 197, 112, 126,
			108, 90, 254, 68, 12, 23, 198, 26, 242, 39, 102, 184, 56, 214,
			146, 63, 9, 195, 165, 49, 71, 95, 80, 48, 32, 119, 230, 5,
			197, 253, 225, 5, 5, 179, 75, 116, 85, 179, 243, 26, 170, 58,
			211, 215, 114, 172, 48, 23, 237, 188, 77, 126, 94, 67, 172, 98,
			240, 243, 26, 202, 25, 252, 188, 150, 47, 26, 252, 188, 150, 36,
			225, 45, 134, 235, 168, 105, 48, 240, 58, 170, 233, 148, 170, 149,
			129, 74, 147, 129, 215, 147, 116, 43, 68, 21, 245, 241, 6, 205,
			3, 138, 157, 98, 164, 129, 28, 28, 215, 217, 242, 186, 90, 157,
			89, 146, 157, 55, 201, 130, 166, 220, 25, 40, 85, 12, 58, 222,
			100, 220, 160, 227, 205, 185, 249, 33, 29, 111, 145, 21, 131, 127,
			183, 72, 213, 224, 223, 45, 21, 153, 196, 252, 187, 181, 180, 172,
			239, 218, 167, 82, 183, 135, 119, 237, 83, 118, 93, 222, 247, 75,
			86, 60, 131, 152, 51, 127, 61, 214, 246, 143, 117, 166, 89, 237,
			109, 106, 112, 231, 25, 52, 213, 208, 252, 56, 13, 16, 182, 193,
			157, 103, 212, 245, 72, 204, 157, 103, 202, 21, 186, 167, 185, 51,
			71, 211, 237, 31, 240, 131, 83, 255, 162, 223, 227, 103, 110, 212,
			61, 149, 113, 173, 186, 98, 58, 15, 68, 79, 28, 123, 96, 197,
			106, 232, 55, 198, 19, 135, 33, 231, 230, 104, 134, 209, 225, 93,
			62, 71, 101, 131, 115, 243, 74, 203, 224, 220, 124, 114, 74, 49,
			119, 196, 240, 44, 106, 41, 20, 136, 240, 102, 17, 159, 86, 45,
			33, 44, 159, 77, 230, 32, 41, 114, 182, 166, 75, 192, 108, 27,
			77, 133, 2, 4, 85, 93, 206, 88, 50, 194, 107, 163, 89, 61,
			30, 68, 120, 237, 4, 5, 86, 163, 157, 213, 185, 1, 88, 141,
			118, 173, 78, 55, 37, 10, 97, 120, 1, 205, 58, 11, 60, 113,
			222, 160, 2, 47, 10, 121, 215, 31, 12, 68, 55, 18, 61, 30,
			187, 173, 225, 196, 129, 148, 45, 160, 182, 78, 89, 144, 12, 96,
			20, 117, 201, 98, 120, 161, 52, 169, 75, 152, 225, 133, 25, 174,
			68, 78, 51, 188, 136, 102, 21, 74, 154, 48, 178, 136, 22, 102,
			85, 203, 116, 6, 42, 53, 10, 80, 136, 197, 4, 37, 141, 25,
			94, 76, 80, 50, 12, 47, 161, 13, 133, 146, 33, 140, 44, 161,
			69, 141, 146, 145, 149, 122, 170, 16, 228, 47, 177, 91, 186, 132,
			25, 94, 186, 189, 174, 80, 198, 24, 94, 65, 101, 133, 50, 70,
			24, 89, 65, 75, 27, 170, 229, 88, 6, 42, 51, 186, 4, 212,
			120, 44, 167, 75, 152, 225, 149, 98, 137, 174, 73, 20, 155, 225,
			85, 52, 227, 204, 240, 159, 134, 23, 110, 191, 255, 158, 123, 209,
			82, 200, 93, 126, 116, 113, 226, 14, 188, 63, 18, 1, 252, 74,
			20, 103, 3, 7, 70, 43, 218, 70, 108, 224, 192, 201, 42, 65,
			12, 188, 154, 117, 116, 9, 56, 240, 212, 180, 90, 165, 44, 195,
			107, 200, 105, 47, 240, 189, 152, 13, 118, 193, 197, 198, 126, 239,
			218, 49, 170, 7, 203, 18, 70, 214, 208, 170, 126, 92, 146, 77,
			3, 134, 30, 12, 66, 226, 181, 196, 176, 178, 152, 225, 181, 70,
			75, 167, 132, 54, 224, 176, 208, 41, 161, 13, 187, 65, 15, 116,
			74, 104, 19, 77, 57, 207, 228, 137, 224, 13, 142, 125, 176, 145,
			51, 17, 191, 198, 249, 177, 219, 121, 203, 195, 247, 97, 36, 206,
			214, 41, 191, 8, 99, 225, 244, 73, 235, 13, 78, 148, 21, 113,
			21, 160, 154, 201, 163, 77, 180, 161, 83, 24, 176, 133, 55, 147,
			199, 38, 176, 133, 55, 179, 77, 35, 121, 180, 57, 49, 73, 111,
			235, 228, 209, 61, 228, 56, 60, 25, 43, 185, 189, 146, 7, 119,
			34, 139, 145, 64, 186, 135, 54, 167, 140, 4, 210, 189, 145, 71,
			45, 247, 178, 117, 35, 129, 116, 47, 126, 153, 16, 39, 144, 238,
			171, 155, 212, 56, 129, 116, 31, 221, 115, 140, 4, 210, 253, 4,
			5, 68, 186, 175, 110, 82, 227, 4, 210, 253, 10, 147, 206, 77,
			22, 30, 160, 5, 103, 94, 73, 121, 228, 251, 223, 174, 241, 190,
			123, 20, 70, 50, 50, 94, 227, 112, 176, 250, 107, 252, 183, 191,
			252, 149, 153, 248, 121, 128, 238, 87, 117, 114, 39, 3, 16, 53,
			35, 241, 243, 160, 206, 141, 196, 207, 131, 185, 121, 218, 161, 49,
			207, 199, 91, 104, 218, 217, 135, 51, 26, 142, 214, 72, 4, 3,
			176, 205, 53, 254, 212, 59, 150, 235, 17, 113, 25, 15, 133, 27,
			50, 142, 129, 101, 4, 39, 24, 158, 186, 129, 144, 222, 47, 4,
			175, 11, 174, 55, 81, 31, 108, 249, 45, 244, 96, 65, 13, 72,
			210, 48, 134, 158, 56, 108, 249, 173, 172, 94, 66, 216, 242, 91,
			147, 83, 244, 7, 84, 101, 25, 30, 161, 150, 179, 18, 27, 72,
			207, 23, 225, 224, 183, 191, 252, 139, 136, 159, 92, 184, 129, 59,
			136, 132, 144, 190, 70, 252, 194, 11, 35, 136, 200, 146, 1, 193,
			59, 60, 66, 91, 250, 41, 83, 90, 226, 232, 1, 193, 59, 60,
			202, 106, 93, 128, 119, 120, 212, 104, 38, 212, 233, 223, 170, 244,
			193, 103, 189, 230, 12, 196, 137, 23, 70, 55, 61, 230, 28, 165,
			94, 237, 7, 52, 167, 78, 158, 111, 188, 48, 98, 203, 116, 76,
			211, 143, 248, 181, 101, 113, 148, 26, 117, 116, 117, 251, 140, 150,
			244, 55, 253, 82, 115, 237, 202, 75, 205, 218, 149, 190, 163, 79,
			52, 215, 174, 60, 209, 252, 64, 107, 245, 54, 243, 29, 45, 140,
			84, 140, 74, 106, 125, 68, 210, 143, 190, 197, 108, 215, 105, 21,
			38, 174, 58, 233, 55, 149, 237, 105, 154, 223, 31, 68, 94, 244,
			254, 249, 83, 169, 151, 162, 98, 139, 120, 57, 11, 220, 176, 253,
			67, 90, 212, 245, 74, 158, 98, 194, 39, 101, 139, 143, 143, 122,
			78, 203, 70, 247, 88, 123, 183, 175, 104, 79, 115, 177, 209, 113,
			18, 245, 221, 190, 162, 190, 15, 53, 143, 27, 109, 254, 6, 209,
			124, 199, 48, 10, 246, 152, 22, 247, 2, 225, 70, 154, 33, 134,
			140, 141, 42, 16, 230, 237, 140, 95, 91, 148, 88, 216, 109, 154,
			251, 145, 72, 180, 198, 170, 87, 6, 255, 104, 223, 93, 154, 55,
			85, 206, 28, 205, 170, 175, 175, 195, 7, 49, 30, 211, 226, 79,
			207, 123, 191, 171, 244, 95, 210, 226, 83, 209, 23, 70, 239, 27,
			39, 208, 184, 174, 210, 145, 183, 180, 255, 156, 143, 249, 199, 246,
			239, 63, 255, 248, 208, 91, 218, 241, 225, 91, 218, 9, 243, 45,
			109, 141, 62, 27, 190, 165, 125, 236, 108, 241, 81, 91, 226, 177,
			255, 17, 65, 200, 7, 226, 59, 117, 50, 108, 196, 127, 190, 113,
			143, 52, 131, 184, 250, 166, 182, 70, 141, 55, 181, 241, 17, 160,
			223, 212, 174, 110, 211, 69, 57, 162, 197, 112, 5, 109, 59, 45,
			254, 35, 17, 233, 160, 89, 30, 217, 193, 153, 52, 108, 141, 40,
			239, 236, 50, 21, 93, 66, 12, 87, 216, 180, 46, 97, 134, 43,
			43, 15, 233, 188, 68, 68, 64, 87, 118, 157, 6, 135, 69, 231,
			110, 191, 111, 146, 156, 80, 227, 201, 247, 254, 25, 166, 75, 208,
			169, 186, 168, 75, 192, 106, 190, 120, 66, 219, 18, 15, 51, 220,
			64, 143, 157, 58, 143, 45, 244, 70, 52, 56, 228, 26, 201, 124,
			129, 184, 53, 146, 249, 194, 145, 215, 88, 221, 86, 104, 132, 97,
			7, 125, 233, 212, 121, 108, 177, 55, 162, 193, 41, 229, 36, 104,
			4, 49, 236, 212, 103, 117, 9, 51, 236, 172, 61, 214, 183, 177,
			83, 169, 153, 225, 109, 236, 148, 93, 29, 254, 11, 201, 52, 210,
			60, 41, 69, 160, 68, 13, 58, 55, 61, 252, 127, 18, 139, 225,
			233, 22, 55, 232, 220, 244, 220, 188, 166, 63, 179, 169, 249, 33,
			253, 153, 85, 105, 114, 73, 127, 218, 104, 94, 179, 21, 2, 37,
			243, 169, 113, 91, 129, 199, 188, 166, 221, 154, 49, 120, 77, 91,
			189, 227, 149, 188, 102, 46, 1, 177, 8, 148, 204, 167, 198, 115,
			9, 8, 140, 61, 151, 128, 192, 82, 207, 181, 231, 116, 252, 183,
			152, 90, 25, 198, 127, 139, 118, 125, 248, 212, 120, 9, 77, 26,
			113, 219, 18, 90, 108, 232, 216, 76, 198, 221, 37, 35, 110, 91,
			42, 55, 140, 184, 109, 73, 253, 215, 138, 180, 144, 229, 145, 167,
			198, 203, 104, 105, 210, 136, 202, 150, 71, 162, 178, 229, 145, 167,
			198, 203, 245, 113, 216, 117, 176, 90, 171, 169, 233, 228, 222, 106,
			213, 158, 208, 119, 69, 183, 83, 119, 134, 119, 69, 183, 237, 154,
			12, 137, 228, 93, 209, 6, 106, 57, 251, 252, 237, 235, 167, 175,
			151, 69, 120, 250, 157, 27, 12, 86, 182, 121, 247, 212, 29, 156,
			8, 238, 245, 120, 228, 199, 215, 70, 111, 221, 19, 238, 7, 60,
			132, 192, 214, 143, 78, 69, 32, 227, 160, 101, 245, 206, 189, 239,
			30, 173, 232, 59, 29, 88, 163, 13, 148, 148, 210, 12, 111, 228,
			42, 116, 120, 167, 180, 193, 244, 85, 11, 40, 96, 163, 209, 148,
			50, 166, 25, 217, 76, 221, 183, 146, 123, 146, 77, 245, 15, 74,
			242, 158, 228, 174, 34, 36, 241, 61, 201, 93, 180, 217, 212, 119,
			33, 105, 168, 180, 141, 123, 146, 187, 217, 156, 113, 79, 114, 183,
			88, 26, 222, 147, 220, 83, 234, 77, 171, 160, 247, 110, 217, 184,
			39, 185, 151, 160, 196, 65, 47, 51, 238, 73, 238, 213, 199, 245,
			61, 201, 3, 240, 215, 250, 158, 228, 129, 250, 135, 27, 121, 79,
			242, 80, 237, 0, 121, 79, 130, 31, 162, 228, 98, 36, 195, 240,
			195, 92, 139, 14, 175, 73, 30, 58, 156, 14, 175, 73, 30, 170,
			188, 129, 188, 38, 217, 74, 64, 44, 25, 173, 38, 247, 34, 25,
			134, 183, 18, 16, 144, 112, 43, 1, 1, 9, 183, 230, 230, 117,
			152, 247, 63, 1, 0, 0, 255, 255, 211, 57, 205, 60, 47, 56,
			0, 0},
	)
}

// FileDescriptorSet returns a descriptor set for this proto package, which
// includes all defined services, and all transitive dependencies.
//
// Will not return nil.
//
// Do NOT modify the returned descriptor.
func FileDescriptorSet() *descriptor.FileDescriptorSet {
	// We just need ONE of the service names to look up the FileDescriptorSet.
	ret, err := discovery.GetDescriptorSet("fleet.Configuration")
	if err != nil {
		panic(err)
	}
	return ret
}
