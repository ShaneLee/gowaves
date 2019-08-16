package state

import (
	"io/ioutil"
	"log"
	"math/big"
	"math/rand"
	"os"
	"testing"

	"github.com/mr-tron/base58/base58"
	"github.com/stretchr/testify/assert"
	"github.com/wavesplatform/gowaves/pkg/crypto"
	"github.com/wavesplatform/gowaves/pkg/keyvalue"
	"github.com/wavesplatform/gowaves/pkg/proto"
	"github.com/wavesplatform/gowaves/pkg/settings"
)

const (
	testSeedLen = 75

	testBloomFilterSize                     = 2e6
	testBloomFilterFalsePositiveProbability = 0.01
	testCacheSize                           = 2 * 1024 * 1024

	testPK   = "AfZtLRQxLNYH5iradMkTeuXGe71uAiATVbr8DpXEEQa8"
	testAddr = "3PDdGex1meSUf4Yq5bjPBpyAbx6us9PaLfo"

	issuerSeed    = "5TUPTbbpiM5UmZDhMmzdsKKNgMvyHwZQncKWfJrxk5bc"
	matcherSeed   = "4TUPTbbpiM5UmZDhMmzdsKKNgMvyHwZQncKWfJrxk4bc"
	minerSeed     = "3TUPTbbpiM5UmZDhMmzdsKKNgMvyHwZQncKWfJrxk3bc"
	senderSeed    = "2TUPTbbpiM5UmZDhMmzdsKKNgMvyHwZQncKWfJrxk2bc"
	recipientSeed = "1TUPTbbpiM5UmZDhMmzdsKKNgMvyHwZQncKWfJrxk1bc"

	assetStr  = "B2u2TBpTYHWCuMuKLnbQfLvdLJ3zjgPiy3iMS2TSYugZ"
	assetStr1 = "3gRJoK6f7XUV7fx5jUzHoPwdb9ZdTFjtTPy2HgDinr1N"

	defaultGenSig = "B2u2TBpTYHWCuMuKLnbQfLvdLJ3zjgPiy3iMS2TSYugZ"

	genesisSignature = "FSH8eAAzZNqnG8xgTZtz5xuLqXySsXgAjmFEC25hXMbEufiGjqWPnGCZFt6gLiVLJny16ipxRNAkkzjjhqTjBE2"
)

var (
	blockID0 = genBlockId(1)
	blockID1 = genBlockId(2)
)

type testAddrData struct {
	sk        crypto.SecretKey
	pk        crypto.PublicKey
	addr      proto.Address
	wavesKey  string
	assetKey  string
	assetKey1 string
}

func newTestAddrData(seedStr string, asset, asset1 []byte) (*testAddrData, error) {
	seedBytes, err := base58.Decode(seedStr)
	if err != nil {
		return nil, err
	}
	sk, pk := crypto.GenerateKeyPair(seedBytes)
	addr, err := proto.NewAddressFromPublicKey('W', pk)
	if err != nil {
		return nil, err
	}
	wavesKey := string((&wavesBalanceKey{addr}).bytes())
	assetKey := string((&assetBalanceKey{addr, asset}).bytes())
	assetKey1 := string((&assetBalanceKey{addr, asset1}).bytes())
	return &testAddrData{sk: sk, pk: pk, addr: addr, wavesKey: wavesKey, assetKey: assetKey, assetKey1: assetKey1}, nil
}

type testAssetData struct {
	asset   *proto.OptionalAsset
	assetID []byte
}

func newTestAssetData(assetStr string) (*testAssetData, error) {
	assetID, err := crypto.NewDigestFromBase58(assetStr)
	if err != nil {
		return nil, err
	}
	asset, err := proto.NewOptionalAssetFromString(assetStr)
	if err != nil {
		return nil, err
	}
	return &testAssetData{asset, assetID.Bytes()}, nil
}

type testGlobalVars struct {
	asset0 *testAssetData
	asset1 *testAssetData

	issuerInfo    *testAddrData
	matcherInfo   *testAddrData
	minerInfo     *testAddrData
	senderInfo    *testAddrData
	recipientInfo *testAddrData
}

var testGlobal testGlobalVars

func TestMain(m *testing.M) {
	var err error
	testGlobal.asset0, err = newTestAssetData(assetStr)
	if err != nil {
		log.Fatalf("newTestAssetData(): %v\n", err)
	}
	testGlobal.asset1, err = newTestAssetData(assetStr1)
	if err != nil {
		log.Fatalf("newTestAssetData(): %v\n", err)
	}
	testGlobal.issuerInfo, err = newTestAddrData(issuerSeed, testGlobal.asset0.assetID, testGlobal.asset1.assetID)
	if err != nil {
		log.Fatalf("newTestAddrData(): %v\n", err)
	}
	testGlobal.matcherInfo, err = newTestAddrData(matcherSeed, testGlobal.asset0.assetID, testGlobal.asset1.assetID)
	if err != nil {
		log.Fatalf("newTestAddrData(): %v\n", err)
	}
	testGlobal.minerInfo, err = newTestAddrData(minerSeed, testGlobal.asset0.assetID, testGlobal.asset1.assetID)
	if err != nil {
		log.Fatalf("newTestAddrData(): %v\n", err)
	}
	testGlobal.senderInfo, err = newTestAddrData(senderSeed, testGlobal.asset0.assetID, testGlobal.asset1.assetID)
	if err != nil {
		log.Fatalf("newTestAddrData(): %v\n", err)
	}
	testGlobal.recipientInfo, err = newTestAddrData(recipientSeed, testGlobal.asset0.assetID, testGlobal.asset1.assetID)
	if err != nil {
		log.Fatalf("newTestAddrData(): %v\n", err)
	}
	os.Exit(m.Run())
}

func defaultTestBloomFilterParams() keyvalue.BloomFilterParams {
	return keyvalue.BloomFilterParams{N: testBloomFilterSize, FalsePositiveProbability: testBloomFilterFalsePositiveProbability}
}

func defaultTestCacheParams() keyvalue.CacheParams {
	return keyvalue.CacheParams{Size: testCacheSize}
}

func defaultTestKeyValParams() keyvalue.KeyValParams {
	return keyvalue.KeyValParams{CacheParams: defaultTestCacheParams(), BloomFilterParams: defaultTestBloomFilterParams()}
}

func defaultAssetInfo(reissuable bool) *assetInfo {
	return &assetInfo{
		assetConstInfo: assetConstInfo{
			issuer:      testGlobal.issuerInfo.pk,
			name:        "asset",
			description: "description",
			decimals:    2,
		},
		assetChangeableInfo: assetChangeableInfo{
			quantity:   *big.NewInt(10000000),
			reissuable: reissuable,
		},
	}
}

type testStorageObjects struct {
	db      keyvalue.IterableKeyVal
	dbBatch keyvalue.Batch
	rw      *blockReadWriter
	hs      *historyStorage
	stateDB *stateDB

	entities *blockchainEntitiesStorage
}

func createStorageObjects() (*testStorageObjects, []string, error) {
	res := make([]string, 2)
	dbDir0, err := ioutil.TempDir(os.TempDir(), "dbDir0")
	if err != nil {
		return nil, nil, err
	}
	res[0] = dbDir0
	rwDir, err := ioutil.TempDir(os.TempDir(), "rw_dir")
	if err != nil {
		return nil, res, err
	}
	res[1] = rwDir
	db, err := keyvalue.NewKeyVal(dbDir0, defaultTestKeyValParams())
	if err != nil {
		return nil, res, err
	}
	dbBatch, err := db.NewBatch()
	if err != nil {
		return nil, res, err
	}
	stateDB, err := newStateDB(db, dbBatch)
	if err != nil {
		return nil, res, err
	}
	rw, err := newBlockReadWriter(rwDir, 8, 8, db, dbBatch)
	if err != nil {
		return nil, res, err
	}
	hs, err := newHistoryStorage(db, dbBatch, rw, stateDB)
	if err != nil {
		return nil, res, err
	}
	entities, err := newBlockchainEntitiesStorage(hs, stateDB, settings.MainNetSettings)
	if err != nil {
		return nil, res, err
	}
	return &testStorageObjects{db, dbBatch, rw, hs, stateDB, entities}, res, nil
}

func (s *testStorageObjects) addBlock(t *testing.T, blockID crypto.Signature) {
	err := s.stateDB.addBlock(blockID)
	assert.NoError(t, err, "stateDB.addBlock() failed")
	err = s.rw.startBlock(blockID)
	assert.NoError(t, err, "startBlock() failed")
	err = s.rw.finishBlock(blockID)
	assert.NoError(t, err, "finishBlock() failed")
}

func (s *testStorageObjects) addBlocks(t *testing.T, blocksNum int) {
	ids := genRandBlockIds(t, blocksNum)
	for _, id := range ids {
		s.addBlock(t, id)
	}
	s.flush(t)
}

func (s *testStorageObjects) createAsset(t *testing.T, assetID crypto.Digest) *assetInfo {
	s.addBlock(t, blockID0)
	assetInfo := defaultAssetInfo(true)
	err := s.entities.assets.issueAsset(assetID, assetInfo, blockID0)
	assert.NoError(t, err, "issueAset() failed")
	s.flush(t)
	return assetInfo
}

func (s *testStorageObjects) activateFeature(t *testing.T, featureID int16) {
	s.addBlock(t, blockID0)
	blockNum, err := s.stateDB.blockIdToNum(blockID0)
	assert.NoError(t, err, "blockIdToNum() failed")
	activationReq := &activatedFeaturesRecord{1, blockNum}
	err = s.entities.features.activateFeature(featureID, activationReq)
	assert.NoError(t, err, "activateFeature() failed")
	s.flush(t)
}

func (s *testStorageObjects) activateSponsorship(t *testing.T) {
	s.activateFeature(t, int16(settings.FeeSponsorship))
	windowSize := settings.MainNetSettings.ActivationWindowSize(1)
	s.addBlocks(t, int(windowSize))
}

func (s *testStorageObjects) flush(t *testing.T) {
	err := s.rw.flush()
	assert.NoError(t, err, "rw.flush() failed")
	s.rw.reset()
	err = s.entities.flush(true)
	assert.NoError(t, err, "entities.flush() failed")
	s.entities.reset()
	err = s.stateDB.flush()
	assert.NoError(t, err, "stateDB.flush() failed")
	s.stateDB.reset()
}

func genRandBlockIds(t *testing.T, amount int) []crypto.Signature {
	ids := make([]crypto.Signature, amount)
	for i := 0; i < amount; i++ {
		id := make([]byte, crypto.SignatureSize)
		_, err := rand.Read(id)
		assert.NoError(t, err, "rand.Read() failed")
		blockID, err := crypto.NewSignatureFromBytes(id)
		assert.NoError(t, err, "NewSignatureFromBytes() failed")
		ids[i] = blockID
	}
	return ids
}

func genBlockId(fillWith byte) crypto.Signature {
	var blockID crypto.Signature
	for i := 0; i < crypto.SignatureSize; i++ {
		blockID[i] = fillWith
	}
	return blockID
}

func generateRandomRecipient(t *testing.T) proto.Recipient {
	seed := make([]byte, testSeedLen)
	_, err := rand.Read(seed)
	assert.NoError(t, err, "rand.Read() failed")
	_, pk := crypto.GenerateKeyPair(seed)
	addr, err := proto.NewAddressFromPublicKey('W', pk)
	assert.NoError(t, err, "NewAddressFromPublicKey() failed")
	return proto.NewRecipientFromAddress(addr)
}