// Copyright (C) 2024 Akave
// See LICENSE for copying information.

package sdk

import (
	"context"
	"fmt"
	"io"

	"github.com/ipfs/boxo/blockservice"
	"github.com/ipfs/boxo/blockstore"
	chunker "github.com/ipfs/boxo/chunker"
	"github.com/ipfs/boxo/exchange/offline"
	"github.com/ipfs/boxo/ipld/merkledag"
	"github.com/ipfs/boxo/ipld/unixfs"
	"github.com/ipfs/boxo/ipld/unixfs/importer/balanced"
	"github.com/ipfs/boxo/ipld/unixfs/importer/helpers"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/sync"
	format "github.com/ipfs/go-ipld-format"
	"github.com/zeebo/errs/v2"

	"github.com/akave-ai/akavesdk/private/encryption"
)

var (
	cidBuilder, _ = merkledag.PrefixForCidVersion(1)
)

// DAGRoot is a helper to build a root CID from chunks.
type DAGRoot struct {
	node   *merkledag.ProtoNode
	fsNode *unixfs.FSNode
}

// NewDAGRoot creates a new DAG root node.
func NewDAGRoot() (*DAGRoot, error) {
	node := new(merkledag.ProtoNode)
	if err := node.SetCidBuilder(cidBuilder); err != nil {
		return nil, err
	}
	fsNode := unixfs.NewFSNode(unixfs.TFile)

	return &DAGRoot{
		node:   node,
		fsNode: fsNode,
	}, nil
}

// AddLink adds a chunk node to the root DAG.
func (root *DAGRoot) AddLink(cid cid.Cid, rawDataSize, protoNodeSize uint64) error {
	root.fsNode.AddBlockSize(rawDataSize)
	link := &format.Link{Size: protoNodeSize, Cid: cid}
	if err := root.node.AddRawLink("", link); err != nil {
		return err
	}

	return nil
}

// Build builds the root CID from the chunk nodes.
func (root *DAGRoot) Build() (cid.Cid, error) {
	if len(root.node.Links()) == 0 {
		return cid.Undef, fmt.Errorf("no chunks added")
	}
	if len(root.node.Links()) == 1 {
		return root.node.Links()[0].Cid, nil
	}

	data, err := root.fsNode.GetBytes()
	if err != nil {
		return cid.Undef, err
	}
	root.node.SetData(data)

	return root.node.Cid(), nil
}

// ChunkDAG is a merkledag of data blocks in a chunk.
type ChunkDAG struct {
	CID           cid.Cid
	RawDataSize   uint64 // size of data read from disk.
	ProtoNodeSize uint64 // size of the ProtoNode in the DAG(RawsDataSize + protonode overhead).
	Blocks        []FileBlockUpload
}

// BuildDAG builds the ChunkDAG of a file.
// TODO: after removing normal api remove encKey from this function.
func BuildDAG(ctx context.Context, reader io.Reader, blockSize int64, encKey []byte) (ChunkDAG, error) {
	var splitter chunker.Splitter
	var err error
	if len(encKey) > 0 { // if encryption is enabled
		splitter, err = encryption.NewSplitter(encKey, reader, blockSize)
		if err != nil {
			return ChunkDAG{}, err
		}
	} else {
		splitter = chunker.NewSizeSplitter(io.NopCloser(reader), blockSize)
	}

	batching := sync.MutexWrap(datastore.NewMapDatastore())
	store := blockstore.NewBlockstore(batching)
	dagServ := merkledag.NewDAGService(blockservice.New(store, offline.Exchange(store)))

	params := helpers.DagBuilderParams{
		Maxlinks:   1024,
		RawLeaves:  false, // 1048590 pb encoded ,1048576 - raw; merkledag.ProtoNode merkledag.RawNode
		CidBuilder: cidBuilder,
		Dagserv:    dagServ,
		NoCopy:     false,
	}

	db, err := params.New(splitter)
	if err != nil {
		return ChunkDAG{}, err
	}

	rootNode, err := balanced.Layout(db)
	if err != nil {
		return ChunkDAG{}, err
	}

	var blocks = make([]FileBlockUpload, 0)
	if len(rootNode.Links()) == 0 {
		// Single block case, the root node itself is the block
		node, err := dagServ.Get(ctx, rootNode.Cid())
		if err != nil {
			return ChunkDAG{}, err
		}

		blocks = append(blocks, FileBlockUpload{
			CID:  node.Cid().String(),
			Data: node.RawData(),
		})
	} else {
		for _, l := range rootNode.Links() {
			node, err := dagServ.Get(ctx, l.Cid)
			if err != nil {
				return ChunkDAG{}, err
			}

			blocks = append(blocks, FileBlockUpload{
				CID:  l.Cid.String(),
				Data: node.RawData(),
			})
		}
	}
	rawDataSize, protoNodeSize, err := nodeSizes(rootNode)
	if err != nil {
		return ChunkDAG{}, err
	}

	return ChunkDAG{
		CID:           rootNode.Cid(),
		RawDataSize:   rawDataSize,
		ProtoNodeSize: protoNodeSize,
		Blocks:        blocks,
	}, nil
}

func nodeSizes(node format.Node) (uint64, uint64, error) {
	protoNode, ok := node.(*merkledag.ProtoNode)
	if !ok {
		return 0, 0, fmt.Errorf("given node %s is not a ProtoNode", node.Cid().String())
	}

	rawDataSize, err1 := unixfs.DataSize(protoNode.Data())
	protoNodeSize, err2 := node.Size()
	if err1 != nil || err2 != nil {
		return 0, 0, errs.Combine(err1, err2)
	}

	return rawDataSize, protoNodeSize, nil
}

func blockByCID(blocks []FileBlockUpload, cid string) (FileBlockUpload, bool) {
	for _, v := range blocks {
		if v.CID == cid {
			return v, true
		}
	}
	return FileBlockUpload{}, false
}
