package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/emirpasic/gods/queues/arrayqueue"
)

type Product struct {
	ID    string `json:"id"`
	Price int    `json:"price"`
}

type RecommendationClient struct {
	httpClient http.Client
	sema       chan struct{}
}

func NewRecommendationClient() RecommendationClient {
	return RecommendationClient{httpClient: http.Client{}, sema: make(chan struct{}, 40)}
}

func (rc RecommendationClient) fetchMany(products []Product) map[string]Product {
	mu := sync.Mutex{}
	result := make(map[string]Product, len(products))

	wg := sync.WaitGroup{}
	for _, product := range products {
		product := product
		rc.sema <- struct{}{}
		wg.Add(1)
		go func() {
			recs := rc.fetchRetry(product.ID)
			mu.Lock()
			for _, rec := range recs {
				rec.Price += product.Price
				if exists, ok := result[rec.ID]; ok {
					rec.Price = min(exists.Price, rec.Price)
					result[rec.ID] = rec
				} else {
					result[rec.ID] = rec
				}
			}
			mu.Unlock()
			wg.Done()
			<-rc.sema
		}()
	}
	wg.Wait()

	return result
}

func (rc RecommendationClient) fetchRetry(productID string) []Product {
	for retry := 0; ; retry++ {
		res, err := rc.fetch(productID)
		if err == nil {
			return res
		}

		fmt.Printf("retrying (%d), product:%s, err:%v\n", retry, productID, err)
		time.Sleep(200 * time.Millisecond)
	}
}

func (rc RecommendationClient) fetch(productID string) ([]Product, error) {
	body := fmt.Sprintf(`
	query {
		getRecommendations(productId: "%s") {
			id
			name
			description
			price
		}
	}
`, productID)

	reqdata := map[string]string{"query": body}
	data, _ := json.Marshal(reqdata)

	req, err := rc.httpClient.Post(
		"http://ppc-ctf.o3.ru:9797/graphql", "application/json", strings.NewReader(string(data)))
	if err != nil {
		return nil, fmt.Errorf("failed to make request for p_id=%s: %w", productID, err)
	}
	defer req.Body.Close()

	if req.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to make request for p_id=%s; status=%s", productID, req.Status)
	}

	type response struct {
		Data struct {
			GetRecommendations []Product `json:"getRecommendations,omitempty"`
		} `json:"data,omitempty"`
	}
	parsedResponse := response{}
	if err = json.NewDecoder(req.Body).Decode(&parsedResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response for p_id=%s: %w", productID, err)
	}

	return parsedResponse.Data.GetRecommendations, nil
}

func lastGraphLayer(rc RecommendationClient, product Product) []Product {
	q, layer := arrayqueue.New(), 0
	q.Enqueue(product)

	for {
		qsize := q.Size()
		layer++
		fmt.Printf("layer: %d, qsize: %d\n", layer, qsize)

		currentLayer := layerFromQ(q, qsize)
		childrenByCurrent := rc.fetchMany(currentLayer)
		if len(childrenByCurrent) == 0 {
			return currentLayer
		}

		for _, child := range childrenByCurrent {
			q.Enqueue(child)
		}
	}
}

func layerFromQ(q *arrayqueue.Queue, qsize int) []Product {
	currentLayer := make([]Product, qsize)
	for i := 0; i < qsize; i++ {
		node, _ := q.Dequeue()
		currentLayer[i] = node.(Product)
	}
	return currentLayer
}

func main() {
	firstProduct := Product{"1337", 1362211}
	rc := NewRecommendationClient()

	lastLayer := lastGraphLayer(rc, firstProduct)
	fmt.Printf("last layer: %d\n", len(lastLayer))

	minCost := math.MaxInt
	for _, e := range lastLayer {
		minCost = min(minCost, e.Price)
	}

	fmt.Printf("min cost: %d\n", minCost)
}
