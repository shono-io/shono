package api

import (
	"encoding/json"
	"fmt"
	"github.com/rs/xid"
	"github.com/shono-io/shono/commons"
	"github.com/shono-io/shono/graph"
	"github.com/shono-io/shono/systems"
	"net/http"
	"strconv"
	"strings"
)

func NewHandler(env graph.Environment, cfg graph.RequestConfig) (_ http.Handler, err error) {
	var store *graph.Store
	var sc graph.StorageClient
	store, sc, err = storageClientByKey(env, cfg.StoreKey)

	bb, err := env.GetBackbone()
	if err != nil {
		return nil, err
	}

	bbc, err := bb.GetClient()
	if err != nil {
		return nil, err
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		parts := strings.Split(r.URL.Path, "/")

		switch r.Method {
		case http.MethodGet:
			if sc == nil {
				// -- if there is no client, we cannot retrieve anything, and we consider this to be intentional
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}

			if len(parts)%2 == 0 {
				// -- if there is an even number of parts, we are retrieving a concept from the storage
				key, err := commons.Parse(parts...)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(err.Error()))
					return
				}

				response, err := sc.Get(r.Context(), store.Collection(), key)
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

						i, err := strconv.Atoi(v[0])
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
						i, err := strconv.Atoi(v[0])
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

				response := map[string]any{
					"total": cur.Count(),
					"items": items,
				}

				if err := json.NewEncoder(w).Encode(response); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(err.Error()))
					return
				}
			}
		case http.MethodPost:
			if cfg.PostEventKey == nil {
				// -- if there is no event key, we cannot publish an event to the backbone, and we consider this to be intentional
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}

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

			if err := bbc.Produce(r.Context(), cfg.PostEventKey, key, body); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}

			w.WriteHeader(http.StatusCreated)

		case http.MethodPut:
			if cfg.PutEventKey == nil {
				// -- if there is no event key, we cannot publish an event to the backbone, and we consider this to be intentional
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}

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

			if err := bbc.Produce(r.Context(), cfg.PutEventKey, key, body); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}

			w.WriteHeader(http.StatusAccepted)

		case http.MethodDelete:
			if cfg.DeleteEventKey == nil {
				// -- if there is no event key, we cannot publish an event to the backbone, and we consider this to be intentional
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}

			key, err := commons.Parse(parts...)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(fmt.Errorf("unable to determine key for url %q: %w", r.URL.String(), err).Error()))
				return
			}

			if err := bbc.Produce(r.Context(), cfg.DeleteEventKey, key, nil); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}

			w.WriteHeader(http.StatusAccepted)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}

		w.WriteHeader(http.StatusNotImplemented)
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
