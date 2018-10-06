package proto

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

const (
	headerLength = 17
	headerMagic  = 0x12345678
)

const (
	AddrSchemeTestnet = 0x54
	ArrcSchemeMainnet = 0x54
)

const (
	GenesisTransactionType = iota + 1
	PaymentTransactionType
	IssueTransactionType
	TransferTransactionType
	ReissueTransactionType
	BurnTransactionType
	ExchangeTransactionType
	LeaseTransactionType
	LeaseCancelTransactionType
	CreateAliasTransactionType
	MassTransferTransactionType
	DataTransactionType
	SetScriptTransactionType
	SponsorFeeTransactionType
)
const (
	contentIDGetPeers      = 0x1
	contentIDPeers         = 0x2
	contentIDGetSignatures = 0x14
	contentIDSignatures    = 0x15
	contentIDGetBlock      = 0x16
	contentIDBlock         = 0x17
	contentIDScore         = 0x18
	contentIDTransaction   = 0x19
	contentIDCheckpoint    = 0x64
)

type Address struct {
	Version       uint8
	AddrScheme    uint8
	PublicKeyHash [20]byte
	CheckSum      [4]byte
}

type Alias struct {
	Version       uint8
	AddrScheme    uint8
	AliasBytesLen uint16
	Alias         []byte
}

type Proof struct {
	Size  uint16
	Proof []byte
}

type BlockSignature [64]byte

type Block struct {
	Version                 uint8
	Timestamp               uint64
	ParentBlockSignature    BlockSignature
	ConsensusBlockLength    uint32
	BaseTarget              uint64
	GenerationSignature     [32]byte
	TransactionsBlockLength uint32
	//Transactions            []Transaction
	BlockSignature BlockSignature
}

type GenesisTransaction struct {
	Type      uint8
	Timestamp uint64
	Amount    uint64
	Recepient [26]byte
}

type IssueTransaction struct {
	Type           uint8
	Signature      BlockSignature
	Type2          uint8
	SenderKey      [32]byte
	NameLength     uint16
	NameBytes      []byte
	DescrLength    uint16
	DescrBytes     []byte
	Quantity       uint64
	Decimals       uint8
	ReissuableFlag uint8
	Fee            uint64
	Timestamp      uint64
}

type ReissueTransaction struct {
	Type           uint8
	Signature      BlockSignature
	Type2          uint8
	SenderKey      [32]byte
	AssetID        [32]byte
	Quantity       uint64
	ReissuableFlag uint8
	Fee            uint64
	Timestamp      uint64
}

type TransferTransaction struct {
	Type                    uint8
	Signature               BlockSignature
	Type2                   uint8
	SenderKey               [32]byte
	AssetFlag               uint8
	AssetID                 [32]byte
	FeeAssetFlag            uint8
	FeeAssetID              [32]byte
	Timestamp               uint64
	Amount                  uint64
	Fee                     uint64
	RecepientAddressOrAlias []byte
	AttachmentLength        uint16
	Attachment              []byte
}

type VersionedTransferTransaction struct {
	Reserved                uint8
	Type                    uint8
	Version                 uint8
	SenderKey               [32]byte
	AssetFlag               uint8
	AssetID                 [32]byte
	Timestamp               uint64
	Amount                  uint64
	Fee                     uint64
	RecepientAddressOrAlias []byte
	AttachmentLength        uint16
	AttachmentBytes         []byte
	ProofVersion            uint8
	ProofNumber             uint16
	Proofs                  []byte
}

type BurnTransaction struct {
	Type      uint8
	SenderKey [32]byte
	AssetID   [32]byte
	Amount    uint64
	Fee       uint64
	Timestamp uint64
	Signature BlockSignature
}

type ExchangeTransaction struct {
	Type                  uint8
	BuyOrderObjectLength  uint32
	SellOrderObjectLength uint32
	BuyOrderObjectBytes   []byte
	SellOrderObjectBytes  []byte
	Price                 uint64
	Amount                uint64
	BuyMatcherFee         uint64
	SellMatcherFee        uint64
	Fee                   uint64
	Timestamp             uint64
	Signature             BlockSignature
}

type LeaseTransaction struct {
	Type                    uint8
	SenderKey               [32]byte
	RecepientAddressOrAlias []byte
	Amount                  uint64
	Fee                     uint64
	Timestamp               uint64
	Signature               BlockSignature
}

type LeaseCancelTransaction struct {
	Version   uint8
	ChainByte uint8
	LeaseId   uint8
	Fee       uint64
	SenderKey [32]byte
	Timestamp uint64
}

type CreateAliasTransaction struct {
	Type             uint8
	SenderKey        [32]byte
	AliasObjectLen   uint16
	AliasObjectBytes []byte
	Fee              uint64
	Timestamp        uint64
	Signature        BlockSignature
}

type MassTransferTransaction struct {
	Type              uint8
	Version           uint8
	SenderKey         [32]byte
	AssetFlag         uint8
	AssetId           [32]byte
	NumberOfTransfers uint16
	Transfers         []byte
	Timestamp         uint8
	Fee               uint8
	AttachmenetLen    uint16
	AttachmentBytes   []byte
	ProofsVersion     uint8
	ProofCount        uint16
	Proofs            []byte
}

type DataEntry struct {
	Key1  string
	Value []byte
}

type DataTransaction struct {
	Reserved       uint8
	Type           uint8
	Version        uint8
	SenderKey      [32]byte
	NumDataEntries uint16

	DataEntries   []DataEntry
	Timestamp     uint64
	Fee           uint64
	ProofsVErsion uint8
	ProofCount    uint8
	SignatureLen  uint16
	Signature     BlockSignature
}

type SponsoredFeeTransaction struct {
	Type               uint8
	Version            uint8
	SenderKey          [32]byte
	AssetID            [32]byte
	MinimalFeeInAssets uint64
	Fee                uint64
	Timestamp          uint64
	Proofs             [64]byte
}

type SetScriptTransaction struct {
	Type              uint8
	Version           uint8
	ChainId           uint8
	SenderKey         [32]byte
	ScriptNotNull     uint8
	ScriptObjectLen   uint16
	ScriptObjectBytes []byte
	Fee               uint64
	Timestamp         uint64
}

type Order struct {
}

type header struct {
	Length        uint32
	Magic         uint32
	ContentID     uint8
	PayloadLength uint32
	PayloadCsum   uint32
}

func (h *header) MarshalBinary() ([]byte, error) {
	data := make([]byte, 17)

	binary.BigEndian.PutUint32(data[0:4], h.Length)
	binary.BigEndian.PutUint32(data[4:8], headerMagic)
	data[8] = h.ContentID
	binary.BigEndian.PutUint32(data[9:13], h.PayloadLength)
	binary.BigEndian.PutUint32(data[13:17], h.PayloadCsum)

	return data, nil
}

func (h *header) UnmarshalBinary(data []byte) error {
	h.Length = binary.BigEndian.Uint32(data[0:4])
	h.Magic = binary.BigEndian.Uint32(data[4:8])
	if h.Magic != headerMagic {
		return fmt.Errorf("received wrong magic: want %x, have %x", headerMagic, h.Magic)
	}
	h.ContentID = data[8]
	h.PayloadLength = binary.BigEndian.Uint32(data[9:13])
	h.PayloadCsum = binary.BigEndian.Uint32(data[13:17])

	return nil
}

type Handshake struct {
	NameLength         uint8
	Name               string
	VersionMajor       uint32
	VersionMinor       uint32
	VersionPatch       uint32
	NodeNameLength     uint8
	NodeName           string
	NodeNonce          uint64
	DeclaredAddrLength uint32
	DeclaredAddrBytes  []byte
	Timestamp          uint64
}

type GetPeersMessage struct{}

func (m *GetPeersMessage) MarshalBinary() ([]byte, error) {
	var header header

	header.Length = headerLength
	header.Magic = headerMagic
	header.ContentID = contentIDGetPeers
	header.PayloadLength = 0
	header.PayloadCsum = 0

	res, err := header.MarshalBinary()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (m *GetPeersMessage) UnmarshalBinary(b []byte) error {
	var header header

	err := header.UnmarshalBinary(b)
	if err != nil {
		return err
	}

	if header.Length != headerLength {
		return fmt.Errorf("getpeers message length is unexpected: want %v have %v", headerLength, header.Length)
	}
	if header.Magic != headerMagic {
		return fmt.Errorf("getpeers message magic is unexpected: want %x have %x", headerMagic, header.Magic)
	}
	if header.ContentID != contentIDGetPeers {
		return fmt.Errorf("getpeers message contentid is unexpected: want %x have %x", contentIDGetPeers, header.ContentID)
	}
	if header.PayloadLength != 0 {
		return fmt.Errorf("getpeers message length is not zero")
	}

	return nil
}

type PeerInfo struct {
	addr net.IP
	port uint16
}

func (m *PeerInfo) MarshalBinary() ([]byte, error) {
	buffer := make([]byte, 6)

	copy(buffer[0:4], m.addr.To4())
	binary.BigEndian.PutUint16(buffer[4:6], m.port)

	return buffer, nil
}

func (m *PeerInfo) UnmarshalBinary(data []byte) error {
	if len(data) < 6 {
		return errors.New("too short")
	}

	m.addr = net.IPv4(data[0], data[1], data[2], data[3])
	m.port = binary.BigEndian.Uint16(data[4:6])

	return nil
}

type PeersMessage struct {
	PeersCount uint32
	Peers      []PeerInfo
}

func (m *PeersMessage) MarshalBinary() ([]byte, error) {
	var h header
	body := make([]byte, 4)

	binary.BigEndian.PutUint32(body[0:4], m.PeersCount)

	for _, k := range m.Peers {
		peer, err := k.MarshalBinary()
		if err != nil {
			return nil, err
		}
		body = append(body, peer...)
	}

	h.Length = headerLength + uint32(len(body))
	h.Magic = headerMagic
	h.ContentID = contentIDPeers
	h.PayloadLength = uint32(len(body))
	h.PayloadCsum = 0

	hdr, err := h.MarshalBinary()
	if err != nil {
		return nil, err
	}

	hdr = append(hdr, body...)

	return hdr, nil
}

func (m *PeersMessage) UnmarshalBinary(data []byte) error {
	var header header
	if err := header.UnmarshalBinary(data); err != nil {
		return err
	}
	data = data[headerLength:]
	if len(data) < 4 {
		return errors.New("peers message has insufficient length")
	}
	peersCount := binary.BigEndian.Uint32(data[0:4])
	data = data[4:]
	for i := uint32(0); i < peersCount; i += 6 {
		var peer PeerInfo
		if err := peer.UnmarshalBinary(data[i : i+6]); err != nil {
			return err
		}
		m.Peers = append(m.Peers, peer)
	}

	return nil
}

type BlockID [64]byte

type GetSignaturesMessage struct {
	Blocks []BlockID
}

func (m *GetSignaturesMessage) MarshalBinary() ([]byte, error) {
	body := make([]byte, 4, 4+len(m.Blocks)*64)
	binary.BigEndian.PutUint32(body[0:4], uint32(len(m.Blocks)))
	for _, b := range m.Blocks {
		body = append(body, b[:]...)
	}

	var h header
	h.Length = headerLength + uint32(len(body))
	h.Magic = headerMagic
	h.ContentID = contentIDGetSignatures
	h.PayloadLength = uint32(len(body))
	h.PayloadCsum = 0

	hdr, err := h.MarshalBinary()
	if err != nil {
		return nil, err
	}

	body = append(hdr, body...)

	return body, nil
}

func (m *GetSignaturesMessage) UnmarshalBinary(data []byte) error {
	var h header
	if err := h.UnmarshalBinary(data); err != nil {
		return err
	}
	if h.Magic != headerMagic {
		return fmt.Errorf("wrong magic in header: %x", h.Magic)
	}
	if h.ContentID != contentIDGetSignatures {
		return fmt.Errorf("wrong content id in header: %x", h.ContentID)
	}
	data = data[17:]
	if len(data) < 4 {
		return fmt.Errorf("message too short %v", len(data))
	}
	blockCount := binary.BigEndian.Uint32(data[0:4])
	data = data[4:]

	for i := uint32(0); i < blockCount; i++ {
		var b BlockID
		if len(data[i:]) < 64 {
			return fmt.Errorf("message too short %v", len(data))
		}
		copy(b[:], data[i:i+64])
		m.Blocks = append(m.Blocks, b)
	}

	return nil
}

type SignaturesMessage struct {
	Signatures []BlockSignature
}

func (m *SignaturesMessage) MarshalBinary() ([]byte, error) {
	body := make([]byte, 4, 4+len(m.Signatures))
	binary.BigEndian.PutUint32(body[0:4], uint32(len(m.Signatures)))
	for _, b := range m.Signatures {
		body = append(body, b[:]...)
	}

	var h header
	h.Length = headerLength + uint32(len(body))
	h.Magic = headerMagic
	h.ContentID = contentIDSignatures
	h.PayloadLength = uint32(len(body))
	h.PayloadCsum = 0

	hdr, err := h.MarshalBinary()
	if err != nil {
		return nil, err
	}

	body = append(hdr, body...)

	return body, nil
}

func (m *SignaturesMessage) UnmarshalBinary(data []byte) error {
	var h header

	if err := h.UnmarshalBinary(data); err != nil {
		return err
	}
	if h.Magic != headerMagic {
		return fmt.Errorf("wrong magic in header: %x", h.Magic)
	}
	if h.ContentID != contentIDSignatures {
		return fmt.Errorf("wrong content id in header: %x", h.ContentID)
	}
	data = data[17:]
	if len(data) < 4 {
		return fmt.Errorf("message too short %v", len(data))
	}
	sigCount := binary.BigEndian.Uint32(data[0:4])
	data = data[4:]

	for i := uint32(0); i < sigCount; i++ {
		var sig BlockSignature
		if len(data[i:]) < 64 {
			return fmt.Errorf("message too short: %v", len(data))
		}
		copy(sig[:], data[i:i+64])
		m.Signatures = append(m.Signatures, sig)
	}

	return nil
}

type GetBlockMessage struct {
	BlockID BlockID
}

func (m *GetBlockMessage) MarshalBinary() ([]byte, error) {
	body := make([]byte, 0, 64)
	body = append(body, m.BlockID[:]...)

	var h header
	h.Length = headerLength + uint32(len(body))
	h.Magic = headerMagic
	h.ContentID = contentIDGetBlock
	h.PayloadLength = uint32(len(body))
	h.PayloadCsum = 0

	hdr, err := h.MarshalBinary()
	if err != nil {
		return nil, err
	}
	body = append(hdr, body...)
	return body, nil
}

func (m *GetBlockMessage) UnmarshalBinary(data []byte) error {
	var h header
	if err := h.UnmarshalBinary(data); err != nil {
		return err
	}

	if h.Magic != headerMagic {
		return fmt.Errorf("wrong magic in header: %x", h.Magic)
	}
	if h.ContentID != contentIDGetBlock {
		return fmt.Errorf("wrong content id in header: %x", h.ContentID)
	}
	data = data[17:]
	if len(data) < 64 {
		return fmt.Errorf("message too short %v", len(data))
	}

	copy(m.BlockID[:], data[:64])

	return nil
}

type BlockMessage struct {
	BlockBytes []byte
}

func (m *BlockMessage) MarshalBinary() ([]byte, error) {
	var h header
	h.Length = headerLength + uint32(len(m.BlockBytes))
	h.Magic = headerMagic
	h.ContentID = contentIDBlock
	h.PayloadLength = uint32(len(m.BlockBytes))
	h.PayloadCsum = 0

	hdr, err := h.MarshalBinary()
	if err != nil {
		return nil, err
	}
	hdr = append(hdr, m.BlockBytes...)
	return hdr, nil
}

func (m *BlockMessage) UnmarshalBinary(data []byte) error {
	var h header
	if err := h.UnmarshalBinary(data); err != nil {
		return err
	}
	if h.Magic != headerMagic {
		return fmt.Errorf("wrong magic in header: %x", h.Magic)
	}
	if h.ContentID != contentIDBlock {
		return fmt.Errorf("wrong content id in header: %x", h.ContentID)
	}

	m.BlockBytes = data[17:]

	return nil
}

type ScoreMessage struct {
	Score []byte
}

func (m *ScoreMessage) MarshalBinary() ([]byte, error) {
	var h header
	h.Length = headerLength + uint32(len(m.Score))
	h.Magic = headerMagic
	h.ContentID = contentIDScore
	h.PayloadLength = uint32(len(m.Score))
	h.PayloadCsum = 0

	hdr, err := h.MarshalBinary()
	if err != nil {
		return nil, err
	}
	hdr = append(hdr, m.Score...)
	return hdr, nil
}

func (m *ScoreMessage) UnmarshalBinary(data []byte) error {
	var h header
	if err := h.UnmarshalBinary(data); err != nil {
		return err
	}
	if h.Magic != headerMagic {
		return fmt.Errorf("wrong magic in header: %x", h.Magic)
	}
	if h.ContentID != contentIDScore {
		return fmt.Errorf("wrong content id in header: %x", h.ContentID)
	}

	m.Score = data[17:]

	return nil
}

type TransactionMessage struct {
	Transaction []byte
}

func (m *TransactionMessage) MarshalBinary() ([]byte, error) {
	var h header
	h.Length = headerLength + uint32(len(m.Transaction))
	h.Magic = headerMagic
	h.ContentID = contentIDTransaction
	h.PayloadLength = uint32(len(m.Transaction))
	h.PayloadCsum = 0

	hdr, err := h.MarshalBinary()
	if err != nil {
		return nil, err
	}
	hdr = append(hdr, m.Transaction...)
	return hdr, nil
}

func (m *TransactionMessage) UnmarshalBinary(data []byte) error {
	var h header
	if err := h.UnmarshalBinary(data); err != nil {
		return err
	}
	if h.Magic != headerMagic {
		return fmt.Errorf("wrong magic in header: %x", h.Magic)
	}
	if h.ContentID != contentIDTransaction {
		return fmt.Errorf("wrong content id in header: %x", h.ContentID)
	}

	m.Transaction = data[17:]

	return nil
}

type CheckpointItem struct {
	Height    uint64
	Signature BlockSignature
}

type CheckPointMessage struct {
	Checkpoints []CheckpointItem
}

func (m *CheckPointMessage) MarshalBinary() ([]byte, error) {
	body := make([]byte, 4, 4+len(m.Checkpoints)*72+100)

	binary.BigEndian.PutUint32(body[0:4], uint32(len(m.Checkpoints)))
	for _, c := range m.Checkpoints {
		var height [8]byte
		binary.BigEndian.PutUint64(height[0:8], c.Height)
		body = append(body, height[:]...)
		body = append(body, c.Signature[:]...)
	}

	var h header
	h.Length = headerLength + uint32(len(body))
	h.Magic = headerMagic
	h.ContentID = contentIDCheckpoint
	h.PayloadLength = uint32(len(body))
	h.PayloadCsum = 0

	hdr, err := h.MarshalBinary()
	if err != nil {
		return nil, err
	}

	hdr = append(hdr, body...)
	return hdr, nil

}

func (m *CheckPointMessage) UnmarshalBinary(data []byte) error {
	var h header
	if err := h.UnmarshalBinary(data); err != nil {
		return err
	}
	if h.Magic != headerMagic {
		return fmt.Errorf("ckeckpoint message magic is unexpected: %x", headerMagic)
	}
	if h.ContentID != contentIDCheckpoint {
		return fmt.Errorf("checkpoint message contentid is unexpected %x", h.ContentID)
	}
	data = data[17:]
	if len(data) < 4 {
		return fmt.Errorf("checkpoint message data too short: %d", len(data))
	}
	checkpointsCount := binary.BigEndian.Uint32(data[0:4])
	data = data[4:]
	for i := uint32(0); i < checkpointsCount; i++ {
		if len(data) < 72 {
			return fmt.Errorf("checkpoint message data too short")
		}
		var ci CheckpointItem
		ci.Height = binary.BigEndian.Uint64(data[0:8])
		copy(ci.Signature[:], data[8:72])
		data = data[72:]
		m.Checkpoints = append(m.Checkpoints, ci)
	}

	return nil
}

func (h *Handshake) marshalBinaryName() ([]byte, error) {
	data := make([]byte, h.NameLength+1)
	data[0] = h.NameLength
	copy(data[1:1+h.NameLength], h.Name)

	return data, nil
}

func (h *Handshake) marshalBinaryVersion() ([]byte, error) {
	data := make([]byte, 12)

	binary.BigEndian.PutUint32(data[0:4], h.VersionMajor)
	binary.BigEndian.PutUint32(data[4:8], h.VersionMinor)
	binary.BigEndian.PutUint32(data[8:12], h.VersionPatch)

	return data, nil
}

func (h *Handshake) marshalBinaryNodeName() ([]byte, error) {
	data := make([]byte, h.NodeNameLength+1)

	data[0] = h.NodeNameLength
	copy(data[1:1+h.NodeNameLength], h.NodeName)

	return data, nil
}

func (h *Handshake) marshalBinaryAddr() ([]byte, error) {
	data := make([]byte, 20+h.DeclaredAddrLength)

	binary.BigEndian.PutUint64(data[0:8], h.NodeNonce)
	binary.BigEndian.PutUint32(data[8:12], h.DeclaredAddrLength)

	copy(data[12:12+h.DeclaredAddrLength], h.DeclaredAddrBytes)
	binary.BigEndian.PutUint64(data[12+h.DeclaredAddrLength:20+h.DeclaredAddrLength], h.Timestamp)

	return data, nil
}

func (h *Handshake) MarshalBinary() ([]byte, error) {
	data1, err := h.marshalBinaryName()
	if err != nil {
		return nil, err
	}
	data2, err := h.marshalBinaryVersion()
	if err != nil {
		return nil, err
	}
	data3, err := h.marshalBinaryNodeName()
	if err != nil {
		return nil, err
	}
	data4, err := h.marshalBinaryAddr()
	if err != nil {
		return nil, err
	}

	data1 = append(data1, data2...)
	data1 = append(data1, data3...)
	data1 = append(data1, data4...)
	return data1, nil
}

func (h *Handshake) UnmarshalBinary(data []byte) error {
	return errors.New("ERR")
}
