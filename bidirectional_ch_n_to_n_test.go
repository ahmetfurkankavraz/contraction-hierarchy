package ch

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
)

func TestManyToManyShortestPath(t *testing.T) {
	g := Graph{}
	err := graphFromCSV(&g, "./data/pgrouting_osm.csv")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("Please wait until contraction hierarchy is prepared")
	g.PrepareContractionHierarchies()
	t.Log("TestShortestPath is starting...")
	u := []int64{106600, 106600, 69618}
	v := []int64{5924, 81611, 69618, 68427, 68490}
	correctAns := [][]float64{
		{61089.42195558673, 94961.78959757874, 78692.8292369651, 61212.00481622628, 71101.1080090782},
		{61089.42195558673, 94961.78959757874, 78692.8292369651, 61212.00481622628, 71101.1080090782},
		{19135.6581215226, -2, -2, -2, -2},
	}
	correctPath := [][]int{
		{418, 866, 591, 314, 353},
		{418, 866, 591, 314, 353},
		{160, -2, -2, -2, -2},
	}
	ans, path := g.ShortestPathManyToMany(u, v)
	// t.Log("ShortestPathManyToMany returned", ans, path)
	for sourceIdx := range u {
		for targetIdx := range v {
			if correctPath[sourceIdx][targetIdx] != -2 && len(path[sourceIdx][targetIdx]) != correctPath[sourceIdx][targetIdx] {
				t.Errorf("Num of vertices in path should be %d, but got %d", correctPath[sourceIdx][targetIdx], len(path[sourceIdx][targetIdx]))
				return
			}
			if correctAns[sourceIdx][targetIdx] != -2 && Round(ans[sourceIdx][targetIdx], 0.00005) != Round(correctAns[sourceIdx][targetIdx], 0.00005) {
				t.Errorf("Length of path should be %f, but got %f", correctAns[sourceIdx][targetIdx], ans[sourceIdx][targetIdx])
				return
			}
		}
	}

	t.Log("TestShortestPath is Ok!")
}

func BenchmarkShortestPathManyToMany(b *testing.B) {
	b.Log("BenchmarkShortestPathManyToMany is starting...")
	rand.Seed(1337)
	for k := 2.0; k <= 8; k++ {
		n := int(math.Pow(2, k))
		g, err := generateSyntheticGraph(n)
		if err != nil {
			b.Error(err)
			return
		}
		b.ResetTimer()
		b.Run(fmt.Sprintf("%s/%d/vertices-%d-edges-%d-shortcuts-%d", "CH shortest path", n, len(g.Vertices), g.GetEdgesNum(), g.GetShortcutsNum()), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				u := []int64{int64(rand.Intn(len(g.Vertices)))}
				v := []int64{
					int64(rand.Intn(len(g.Vertices))),
					int64(rand.Intn(len(g.Vertices))),
					int64(rand.Intn(len(g.Vertices))),
					int64(rand.Intn(len(g.Vertices))),
					int64(rand.Intn(len(g.Vertices))),
				}
				ans, path := g.ShortestPathManyToMany(u, v)
				_, _ = ans, path
			}
		})
	}
}

func BenchmarkOldWayShortestPathManyToMany(b *testing.B) {
	b.Log("BenchmarkOldWayShortestPathManyToMany is starting...")
	rand.Seed(1337)
	for k := 2.0; k <= 8; k++ {
		n := int(math.Pow(2, k))
		g, err := generateSyntheticGraph(n)
		if err != nil {
			b.Error(err)
			return
		}
		b.ResetTimer()
		b.Run(fmt.Sprintf("%s/%d/vertices-%d-edges-%d-shortcuts-%d", "CH shortest path", n, len(g.Vertices), g.GetEdgesNum(), g.GetShortcutsNum()), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				u := int64(rand.Intn(len(g.Vertices)))
				v := []int64{
					int64(rand.Intn(len(g.Vertices))),
					int64(rand.Intn(len(g.Vertices))),
					int64(rand.Intn(len(g.Vertices))),
					int64(rand.Intn(len(g.Vertices))),
					int64(rand.Intn(len(g.Vertices))),
				}
				for vv := range v {
					ans, path := g.ShortestPath(u, v[vv])
					_, _ = ans, path
				}
			}
		})
	}
}

func BenchmarkStaticCaseShortestPathManyToMany(b *testing.B) {
	g := Graph{}
	err := graphFromCSV(&g, "./data/pgrouting_osm.csv")
	if err != nil {
		b.Error(err)
	}
	b.Log("Please wait until contraction hierarchy is prepared")
	g.PrepareContractionHierarchies()
	b.Log("BenchmarkStaticCaseShortestPathManyToMany is starting...")
	b.ResetTimer()

	b.Run(fmt.Sprintf("%s/vertices-%d", "CH shortest path (many to many)", len(g.Vertices)), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			u := []int64{106600}
			v := []int64{5924, 81611, 69618, 68427, 68490}
			ans, path := g.ShortestPathManyToMany(u, v)
			_, _ = ans, path
		}
	})
}

func BenchmarkStaticCaseOldWayShortestPathManyToMany(b *testing.B) {
	g := Graph{}
	err := graphFromCSV(&g, "data/pgrouting_osm.csv")
	if err != nil {
		b.Error(err)
	}
	b.Log("Please wait until contraction hierarchy is prepared")
	g.PrepareContractionHierarchies()
	b.Log("BenchmarkStaticCaseOldWayShortestPathManyToMany is starting...")
	b.ResetTimer()

	b.Run(fmt.Sprintf("%s/vertices-%d", "CH shortest path (many to many)", len(g.Vertices)), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			u := int64(106600)
			v := []int64{5924, 81611, 69618, 68427, 68490}
			for vv := range v {
				ans, path := g.ShortestPath(u, v[vv])
				_, _ = ans, path
			}
		}
	})
}
