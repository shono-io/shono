package api

import (
	"encoding/json"
	"fmt"
	"github.com/rs/xid"
	"github.com/shono-io/shono/commons"
	"github.com/shono-io/shono/graph"
	"github.com/shono-io/shono/systems"
	"github.com/shono-io/shono/systems/backbone"
	"net/http"
	"strconv"
	"strings"
)

func NewRequestHandler(env graph.Environment, bbc backbone.Client, req graph.Request) (h http.Handler, err error) {
	switch req.Kind {
	case graph.ListOperationType, graph.GetOperationType:
		var store *graph.Store
		var sc graph.StorageClient
		store, sc, err = storageClientByKey(env, req.StoreKey)
		return newGetRequestHandler(store, sc)

	case graph.CreateOperationType:
		return newPostRequestHandler(bbc, req.EventKey)
	case graph.UpdateOperationType:
		return newPutRequestHandler(bbc, req.EventKey)
	case graph.DeleteOperationType:
		return newDeleteRequestHandler(bbc, req.EventKey)
	default:
		return nil, fmt.Errorf("unknown request type: %s", req.Kind)
	}
}

func newGetRequestHandler(store *graph.Store, sc graph.StorageClient) (http.Handler, error) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")

		if len(parts)%2 == 0 {
			// -- if there is an even number of parts, we are retrieving a concept from the storage
			key, err := commons.Parse(parts...)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}

			response, err := sc.Get(r.Context(), store.Collection(), key.String())
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			if response == nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			if err := json.NewEncoder(w).Encode(response); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
		} else {
			// -- if there is an uneven number of parts, we are asking for a list of concepts from the storage
			var paging *graph.PagingOpts
			filters := map[string]any{}
			for k, v := range r.URL.Query() {
				switch k {
				case "_size":
					if paging == nil {
						paging = &graph.PagingOpts{}
					}

					i, err := strconv.ParseInt(v[0], 10, 64)
					if err != nil {
						w.WriteHeader(http.StatusBadRequest)
						w.Write([]byte(fmt.Errorf("invalid _size request parameter: %w", err).Error()))
						return
					}
					paging.Size = i
				case "_offset":
					if paging == nil {
						paging = &graph.PagingOpts{}
					}
					i, err := strconv.ParseInt(v[0], 10, 64)
					if err != nil {
						w.WriteHeader(http.StatusBadRequest)
						w.Write([]byte(fmt.Errorf("invalid _offset request parameter: %w", err).Error()))
						return
					}
					paging.Offset = i
				default:
					filters[k] = v
				}
			}

			cur, err := sc.List(r.Context(), store.Collection(), filters, paging)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			defer cur.Close()

			var items []map[string]any
			for cur.HasNext() {
				item, err := cur.Read()
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(err.Error()))
					return
				}
				items = append(items, item)
			}

			if err := json.NewEncoder(w).Encode(items); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
		}
	}), nil
}

func newPostRequestHandler(bbc backbone.Client, eventKey commons.Key) (http.Handler, error) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")
		defer r.Body.Close()

		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		parts = append(parts, xid.New().String())
		key, err := commons.Parse(parts...)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Errorf("unable to determine key for url %q: %w", r.URL.String(), err).Error()))
			return
		}

		if err := bbc.Produce(r.Context(), eventKey, key, body); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusCreated)
	}), nil
}

func newPutRequestHandler(bbc backbone.Client, eventKey commons.Key) (http.Handler, error) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")
		defer r.Body.Close()

		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		key, err := commons.Parse(parts...)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Errorf("unable to determine key for url %q: %w", r.URL.String(), err).Error()))
			return
		}

		if err := bbc.Produce(r.Context(), eventKey, key, body); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusAccepted)
	}), nil
}

func newDeleteRequestHandler(bbc backbone.Client, eventKey commons.Key) (http.Handler, error) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")
		key, err := commons.Parse(parts...)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Errorf("unable to determine key for url %q: %w", r.URL.String(), err).Error()))
			return
		}

		if err := bbc.Produce(r.Context(), eventKey, key, nil); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusAccepted)
	}), nil
}

func storageClientByKey(env graph.Environment, storeKey commons.Key) (*graph.Store, graph.StorageClient, error) {
	// -- get the store from the environment
	store, err := env.GetStore(storeKey)
	if err != nil {
		return nil, nil, err
	}

	// -- get the storage corresponding to the store
	storage, err := env.GetStorage(store.StorageKey())
	if err != nil {
		return store, nil, err
	}

	// -- get the storage client
	ss, fnd := systems.Storage[storage.Kind()]
	if !fnd {
		return store, nil, fmt.Errorf("storage system '%s' not found", storage.Kind())
	}

	// -- get the client to the storage system
	r, err := ss.GetClient(storage.Config())
	return store, r, err
}
