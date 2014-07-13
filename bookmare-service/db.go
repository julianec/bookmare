package main

import (
    "database/cassandra"
    "errors"
    "time"
    "github.com/julianec/bookmare"
)

type BookmarkDB struct {
    db *cassandra.RetryCassandraClient
}

func NewBookmarkDB(host, keyspace string) (*BookmarkDB, error)  {
    var ret *BookmarkDB
    var client *cassandra.RetryCassandraClient
    var err error
    var ire *cassandra.InvalidRequestException

    client, err = cassandra.NewRetryCassandraClient(host)
    if err != nil {
        return nil, err
    }

    ire, err = client.SetKeyspace(keyspace)
    if err != nil {
        return nil, err
    }
    if ire != nil {
        return nil, errors.New(ire.Why)
    }

    ret = &BookmarkDB{
        db: client,
    }
    return ret, err
}

func makeMutation(name string, value []byte, now time.Time) *cassandra.Mutation {
    var mutation *cassandra.Mutation = cassandra.NewMutation()
    mutation.ColumnOrSupercolumn = cassandra.NewColumnOrSuperColumn()
    mutation.ColumnOrSupercolumn.Column = cassandra.NewColumn()
    mutation.ColumnOrSupercolumn.Column.Name = []byte(name)
    mutation.ColumnOrSupercolumn.Column.Value = value
    mutation.ColumnOrSupercolumn.Column.Timestamp = now.UnixNano()
    return mutation
}

func makeMutationStr(name string, value string, now time.Time) *cassandra.Mutation {
    return makeMutation(name, []byte(value), now)
}

func (b *BookmarkDB) SaveBookmark(bookmark *bookmare.Bookmark) error {
    var mutation_map map[string]map[string][]*cassandra.Mutation
    var mlist []*cassandra.Mutation
    var id cassandra.UUID
    var idstr string
    var err error
    var now time.Time = time.Now()
    var ire *cassandra.InvalidRequestException
    var ue *cassandra.UnavailableException
    var te *cassandra.TimedOutException

    // Create row key
    id, err = cassandra.GenTimeUUID(&now)
    if err != nil {
        return err
    }

    idstr = string([]byte(id))
    mutation_map = make(map[string]map[string][]*cassandra.Mutation)
    mutation_map[idstr]=make(map[string][]*cassandra.Mutation)
    mlist = make([]*cassandra.Mutation,0)

    mlist = append(mlist, makeMutationStr("url", *bookmark.Url, now))
    mlist = append(mlist, makeMutationStr("owner", *bookmark.Owner, now))
    mlist = append(mlist, makeMutationStr("title", *bookmark.Title, now))
    mlist = append(mlist, makeMutationStr("description", *bookmark.Description, now))

    mutation_map[idstr]["bookmark"] = mlist
    ire, ue, te, err = b.db.AtomicBatchMutate(mutation_map, cassandra.ConsistencyLevel_QUORUM)

    // error handling
    if ire != nil {
        return errors.New(ire.Why)
    }
    if ue != nil {
        return errors.New("Unavailable")
    }
    if te != nil {
        return errors.New("Timeout")
    }
    return err
}
