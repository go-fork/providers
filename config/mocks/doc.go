// Package mocks cung cấp các implement giả lập (mock) cho các interface trong package config,
// được tạo tự động bởi mockery và sử dụng testify/mock framework để hỗ trợ viết unit test.
//
// # Đối tượng chính
//
//   - MockManager: Mock implementation của interface Manager được tạo bởi mockery v2.53.4,
//     cung cấp khả năng kiểm soát hoàn toàn hành vi của các phương thức thông qua testify/mock.
//
// # Tính năng
//
//   - Mock tự động cho tất cả các phương thức của interface Manager
//   - Hỗ trợ expecter pattern để thiết lập expectations dễ dàng
//   - Xác thực các method calls với mock.AssertExpectations()
//   - Hỗ trợ Return, Run, và RunAndReturn patterns
//   - Panic khi method được gọi mà không có expectation được thiết lập
//
// # Ví dụ sử dụng với Expecter Pattern
//
//	func TestWithMockManager(t *testing.T) {
//		// Tạo mock manager
//		mockCfg := &mocks.MockManager{}
//
//		// Thiết lập expectations sử dụng EXPECT()
//		mockCfg.EXPECT().GetString("app.name").Return("TestApp", true)
//		mockCfg.EXPECT().GetInt("database.port").Return(5432, true)
//		mockCfg.EXPECT().Has("feature.enabled").Return(true)
//
//		// Sử dụng mock trong test
//		name, ok := mockCfg.GetString("app.name")
//		assert.True(t, ok)
//		assert.Equal(t, "TestApp", name)
//
//		port, ok := mockCfg.GetInt("database.port")
//		assert.True(t, ok)
//		assert.Equal(t, 5432, port)
//
//		// Xác thực tất cả expectations đã được gọi
//		mockCfg.AssertExpectations(t)
//	}
//
// # Ví dụ sử dụng với Traditional Mock
//
//	func TestWithTraditionalMock(t *testing.T) {
//		mockCfg := &mocks.MockManager{}
//
//		// Thiết lập mock responses
//		mockCfg.On("GetString", "app.name").Return("TestApp", true)
//		mockCfg.On("Set", "new.key", mock.Anything).Return(nil)
//
//		// Test code sử dụng mock
//		name, _ := mockCfg.GetString("app.name")
//		err := mockCfg.Set("new.key", "value")
//
//		assert.Equal(t, "TestApp", name)
//		assert.NoError(t, err)
//		mockCfg.AssertExpectations(t)
//	}
//
// # Lưu ý quan trọng
//
//   - File manager.go được tạo tự động bởi mockery, KHÔNG CHỈNH SỬA thủ công
//   - Để tái tạo mock khi interface thay đổi, chạy: mockery --name Manager
//   - Mock sẽ panic nếu method được gọi mà không có expectation tương ứng
//   - Luôn gọi AssertExpectations(t) để đảm bảo tất cả expectations đã được thực hiện
package mocks
