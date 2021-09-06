package blockchain

type mockDB struct {
	mockLoadChain func() []byte
	mockFindBlock func(hash string) []byte
}

func (m mockDB) FindBlock(hash string) []byte {
	return m.mockFindBlock(hash)
}
func (m mockDB) LoadBlockchain() []byte {
	return m.mockLoadChain()
}
func (mockDB) SaveBlock(hash string, data []byte) {}
func (mockDB) SaveBlockchain(data []byte)         {}
func (mockDB) EmptyBlocks()                       {}
