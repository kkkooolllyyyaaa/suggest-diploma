package subtree

import (
	"reflect"
	"testing"

	doppelgangerNormalizer "go.avito.ru/gl/doppelganger-normalizer"

	"go.avito.ru/av/service-suggest/internal/config/envs"
	"go.avito.ru/av/service-suggest/internal/infm"
	"go.avito.ru/av/service-suggest/internal/infm/catDict"
	"go.avito.ru/av/service-suggest/internal/infm/dictionary"
	"go.avito.ru/av/service-suggest/internal/infm/hierarchy"
)

var queryStats = map[string][]catDict.NodeStats{
	"диван бу": {
		// Для дома и дачи
		{
			NodeID:               1047737,
			QueryFreq:            943993,
			QueryContactFreq:     23141,
			QuerySearchFreq:      249763,
			McatQueryFreq:        940215,
			McatQueryContactFreq: 23041,
			McatQuerySearchFreq:  248985,
			QueryMcatRate:        0.995998,
			QueryMcatContactRate: 0.995679,
			QueryMcatSearchRate:  0.996885,
		},

		// Мебель и интерьер
		{
			NodeID:               1047738,
			QueryFreq:            943993,
			QueryContactFreq:     23141,
			QuerySearchFreq:      249763,
			McatQueryFreq:        879009,
			McatQueryContactFreq: 23027,
			McatQuerySearchFreq:  188199,
			QueryMcatRate:        0.931161,
			QueryMcatContactRate: 0.995074,
			QueryMcatSearchRate:  0.75351,
		},

		// Кровати, диваны и кресла
		{
			NodeID:               1047739,
			QueryFreq:            943993,
			QueryContactFreq:     23141,
			QuerySearchFreq:      249763,
			McatQueryFreq:        846385,
			McatQueryContactFreq: 22852,
			McatQuerySearchFreq:  160825,
			QueryMcatRate:        0.896601,
			QueryMcatContactRate: 0.987511,
			QueryMcatSearchRate:  0.64391,
		},

		// Диваны
		{
			NodeID:               1090116,
			QueryFreq:            943993,
			QueryContactFreq:     23141,
			QuerySearchFreq:      249763,
			McatQueryFreq:        793658,
			McatQueryContactFreq: 21866,
			McatQuerySearchFreq:  137678,
			QueryMcatRate:        0.840746,
			QueryMcatContactRate: 0.944903,
			QueryMcatSearchRate:  0.551235,
		},

		// Прямые диваны
		{
			NodeID:               1288006,
			QueryFreq:            943993,
			QueryContactFreq:     23141,
			QuerySearchFreq:      249763,
			McatQueryFreq:        603071,
			McatQueryContactFreq: 16704,
			McatQuerySearchFreq:  101951,
			QueryMcatRate:        0.638851,
			QueryMcatContactRate: 0.721836,
			QueryMcatSearchRate:  0.408191,
		},

		// Угловые диваны
		{
			NodeID:               1288007,
			QueryFreq:            943993,
			QueryContactFreq:     23141,
			QuerySearchFreq:      249763,
			McatQueryFreq:        148732,
			McatQueryContactFreq: 4780,
			McatQuerySearchFreq:  5332,
			QueryMcatRate:        0.157556,
			QueryMcatContactRate: 0.20656,
			QueryMcatSearchRate:  0.021348,
		},

		// Кресла
		{
			NodeID:               1210301,
			QueryFreq:            943993,
			QueryContactFreq:     23141,
			QuerySearchFreq:      249763,
			McatQueryFreq:        23848,
			McatQueryContactFreq: 770,
			McatQuerySearchFreq:  748,
			QueryMcatRate:        0.025263,
			QueryMcatContactRate: 0.033274,
			QueryMcatSearchRate:  0.002995,
		},

		// Новые прямые диваны
		{
			NodeID:               2170275,
			QueryFreq:            943993,
			QueryContactFreq:     23141,
			QuerySearchFreq:      249763,
			McatQueryFreq:        16824,
			McatQueryContactFreq: 497,
			McatQuerySearchFreq:  1914,
			QueryMcatRate:        0.017822,
			QueryMcatContactRate: 0.021477,
			QueryMcatSearchRate:  0.007663,
		},
	},
	"бассейн": {
		// Для дома и дачи
		{
			NodeID:               1047737,
			QueryFreq:            78475,
			QueryContactFreq:     1342,
			QuerySearchFreq:      38215,
			McatQueryFreq:        41929,
			McatQueryContactFreq: 637,
			McatQuerySearchFreq:  22819,
			QueryMcatRate:        0.534298,
			QueryMcatContactRate: 0.474665,
			QueryMcatSearchRate:  0.597122,
		},

		// Товары для ремонта
		{
			NodeID:               1047751,
			QueryFreq:            78475,
			QueryContactFreq:     1342,
			QuerySearchFreq:      38215,
			McatQueryFreq:        39824,
			McatQueryContactFreq: 633,
			McatQuerySearchFreq:  20834,
			QueryMcatRate:        0.507474,
			QueryMcatContactRate: 0.471684,
			QueryMcatSearchRate:  0.545179,
		},

		// Сантехника, водоснабжение и сауна
		{
			NodeID:               1047754,
			QueryFreq:            78475,
			QueryContactFreq:     1342,
			QuerySearchFreq:      38215,
			McatQueryFreq:        39366,
			McatQueryContactFreq: 625,
			McatQuerySearchFreq:  20616,
			QueryMcatRate:        0.501637,
			QueryMcatContactRate: 0.465723,
			QueryMcatSearchRate:  0.539474,
		},

		// None
		{
			NodeID:               1288806,
			QueryFreq:            78475,
			QueryContactFreq:     1342,
			QuerySearchFreq:      38215,
			McatQueryFreq:        39275,
			McatQueryContactFreq: 623,
			McatQuerySearchFreq:  20585,
			QueryMcatRate:        0.500478,
			QueryMcatContactRate: 0.464232,
			QueryMcatSearchRate:  0.538663,
		},

		// Бассейны и комплектующие
		{
			NodeID:               2609066,
			QueryFreq:            78475,
			QueryContactFreq:     1342,
			QuerySearchFreq:      38215,
			McatQueryFreq:        38891,
			McatQueryContactFreq: 613,
			McatQuerySearchFreq:  20501,
			QueryMcatRate:        0.495585,
			QueryMcatContactRate: 0.456781,
			QueryMcatSearchRate:  0.536465,
		},

		// Бассейны
		{
			NodeID:               2618230,
			QueryFreq:            78475,
			QueryContactFreq:     1342,
			QuerySearchFreq:      38215,
			McatQueryFreq:        36603,
			McatQueryContactFreq: 605,
			McatQuerySearchFreq:  18453,
			QueryMcatRate:        0.466429,
			QueryMcatContactRate: 0.45082,
			QueryMcatSearchRate:  0.482873,
		},

		// Услуги
		{
			NodeID:               1047505,
			QueryFreq:            78475,
			QueryContactFreq:     1342,
			QuerySearchFreq:      38215,
			McatQueryFreq:        20504,
			McatQueryContactFreq: 477,
			McatQuerySearchFreq:  6194,
			QueryMcatRate:        0.261281,
			QueryMcatContactRate: 0.35544,
			QueryMcatSearchRate:  0.162083,
		},

		// Услуги
		{
			NodeID:               1047506,
			QueryFreq:            78475,
			QueryContactFreq:     1342,
			QuerySearchFreq:      38215,
			McatQueryFreq:        20383,
			McatQueryContactFreq: 477,
			McatQuerySearchFreq:  6073,
			QueryMcatRate:        0.259739,
			QueryMcatContactRate: 0.35544,
			QueryMcatSearchRate:  0.158917,
		},

		// Услуги в сфере красоты и здоровья
		{
			NodeID:               1047554,
			QueryFreq:            78475,
			QueryContactFreq:     1342,
			QuerySearchFreq:      38215,
			McatQueryFreq:        14790,
			McatQueryContactFreq: 332,
			McatQuerySearchFreq:  4830,
			QueryMcatRate:        0.188468,
			QueryMcatContactRate: 0.247392,
			QueryMcatSearchRate:  0.12639,
		},

		// Спа-услуги, массаж
		{
			NodeID:               1047561,
			QueryFreq:            78475,
			QueryContactFreq:     1342,
			QuerySearchFreq:      38215,
			McatQueryFreq:        13949,
			McatQueryContactFreq: 324,
			McatQuerySearchFreq:  4229,
			QueryMcatRate:        0.177751,
			QueryMcatContactRate: 0.241431,
			QueryMcatSearchRate:  0.110663,
		},

		// Недвижимость
		{
			NodeID:               1054803,
			QueryFreq:            78475,
			QueryContactFreq:     1342,
			QuerySearchFreq:      38215,
			McatQueryFreq:        9217,
			McatQueryContactFreq: 83,
			McatQuerySearchFreq:  6727,
			QueryMcatRate:        0.117451,
			QueryMcatContactRate: 0.061848,
			QueryMcatSearchRate:  0.17603,
		},

		// Дома, дачи, коттеджи
		{
			NodeID:               1054804,
			QueryFreq:            78475,
			QueryContactFreq:     1342,
			QuerySearchFreq:      38215,
			McatQueryFreq:        7474,
			McatQueryContactFreq: 56,
			McatQuerySearchFreq:  5794,
			QueryMcatRate:        0.095241,
			QueryMcatContactRate: 0.041729,
			QueryMcatSearchRate:  0.151616,
		},

		// Жёсткие бассейны
		{
			NodeID:               4698124,
			QueryFreq:            78475,
			QueryContactFreq:     1342,
			QuerySearchFreq:      38215,
			McatQueryFreq:        5758,
			McatQueryContactFreq: 190,
			McatQuerySearchFreq:  58,
			QueryMcatRate:        0.073374,
			QueryMcatContactRate: 0.14158,
			QueryMcatSearchRate:  0.001518,
		},

		// Снять дом, дачу, коттедж
		{
			NodeID:               1054823,
			QueryFreq:            78475,
			QueryContactFreq:     1342,
			QuerySearchFreq:      38215,
			McatQueryFreq:        5726,
			McatQueryContactFreq: 50,
			McatQuerySearchFreq:  4226,
			QueryMcatRate:        0.072966,
			QueryMcatContactRate: 0.037258,
			QueryMcatSearchRate:  0.110585,
		},

		// Каркасные бассейны
		{
			NodeID:               4698122,
			QueryFreq:            78475,
			QueryContactFreq:     1342,
			QuerySearchFreq:      38215,
			McatQueryFreq:        4747,
			McatQueryContactFreq: 155,
			McatQuerySearchFreq:  97,
			QueryMcatRate:        0.060491,
			QueryMcatContactRate: 0.115499,
			QueryMcatSearchRate:  0.002538,
		},

		// Услуги по организации мероприятий
		{
			NodeID:               1047535,
			QueryFreq:            78475,
			QueryContactFreq:     1342,
			QuerySearchFreq:      38215,
			McatQueryFreq:        3346,
			McatQueryContactFreq: 98,
			McatQuerySearchFreq:  406,
			QueryMcatRate:        0.042638,
			QueryMcatContactRate: 0.073025,
			QueryMcatSearchRate:  0.010624,
		},

		// Личные вещи
		{
			NodeID:               1054439,
			QueryFreq:            78475,
			QueryContactFreq:     1342,
			QuerySearchFreq:      38215,
			McatQueryFreq:        2344,
			McatQueryContactFreq: 60,
			McatQuerySearchFreq:  544,
			QueryMcatRate:        0.029869,
			QueryMcatContactRate: 0.044709,
			QueryMcatSearchRate:  0.014235,
		},

		// Хобби и отдых
		{
			NodeID:               1054252,
			QueryFreq:            78475,
			QueryContactFreq:     1342,
			QuerySearchFreq:      38215,
			McatQueryFreq:        2294,
			McatQueryContactFreq: 53,
			McatQuerySearchFreq:  704,
			QueryMcatRate:        0.029232,
			QueryMcatContactRate: 0.039493,
			QueryMcatSearchRate:  0.018422,
		},

		// Детские товары, игрушки
		{
			NodeID:               1054789,
			QueryFreq:            78475,
			QueryContactFreq:     1342,
			QuerySearchFreq:      38215,
			McatQueryFreq:        2264,
			McatQueryContactFreq: 60,
			McatQuerySearchFreq:  464,
			QueryMcatRate:        0.02885,
			QueryMcatContactRate: 0.044709,
			QueryMcatSearchRate:  0.012142,
		},

		// Бассейны Bestway
		{
			NodeID:               4698121,
			QueryFreq:            78475,
			QueryContactFreq:     1342,
			QuerySearchFreq:      38215,
			McatQueryFreq:        2024,
			McatQueryContactFreq: 67,
			McatQuerySearchFreq:  14,
			QueryMcatRate:        0.025792,
			QueryMcatContactRate: 0.049925,
			QueryMcatSearchRate:  0.000366,
		},

		// Спорт и отдых
		{
			NodeID:               1054305,
			QueryFreq:            78475,
			QueryContactFreq:     1342,
			QuerySearchFreq:      38215,
			McatQueryFreq:        1775,
			McatQueryContactFreq: 48,
			McatQuerySearchFreq:  335,
			QueryMcatRate:        0.022619,
			QueryMcatContactRate: 0.035768,
			QueryMcatSearchRate:  0.008766,
		},

		// Товары для купания
		{
			NodeID:               1054798,
			QueryFreq:            78475,
			QueryContactFreq:     1342,
			QuerySearchFreq:      38215,
			McatQueryFreq:        1688,
			McatQueryContactFreq: 48,
			McatQuerySearchFreq:  248,
			QueryMcatRate:        0.02151,
			QueryMcatContactRate: 0.035768,
			QueryMcatSearchRate:  0.00649,
		},

		// Организация досуга и отдыха
		{
			NodeID:               4442348,
			QueryFreq:            78475,
			QueryContactFreq:     1342,
			QuerySearchFreq:      38215,
			McatQueryFreq:        1585,
			McatQueryContactFreq: 49,
			McatQuerySearchFreq:  115,
			QueryMcatRate:        0.020198,
			QueryMcatContactRate: 0.036513,
			QueryMcatSearchRate:  0.003009,
		},
	},

	// Терминатор на ноды ниже "Спорт и отдых"
	"query with terminator": {
		// Личные вещи
		{
			NodeID:               1054439,
			QueryFreq:            78475,
			QueryContactFreq:     1342,
			QuerySearchFreq:      38215,
			McatQueryFreq:        2344,
			McatQueryContactFreq: 60,
			McatQuerySearchFreq:  544,
			QueryMcatRate:        0.029869,
			QueryMcatContactRate: 0.044709,
			QueryMcatSearchRate:  0.014235,
		},
		// Хобби и отдых
		{
			NodeID:               1054252,
			QueryFreq:            78475,
			QueryContactFreq:     1342,
			QuerySearchFreq:      38215,
			McatQueryFreq:        2294,
			McatQueryContactFreq: 53,
			McatQuerySearchFreq:  704,
			QueryMcatRate:        0.029232,
			QueryMcatContactRate: 0.039493,
			QueryMcatSearchRate:  0.018422,
		},

		// Хобби и отдых -> Спорт и отдых
		{
			NodeID:               1054305,
			QueryFreq:            78475,
			QueryContactFreq:     1342,
			QuerySearchFreq:      38215,
			McatQueryFreq:        1775,
			McatQueryContactFreq: 48,
			McatQuerySearchFreq:  335,
			QueryMcatRate:        0.022619,
			QueryMcatContactRate: 0.035768,
			QueryMcatSearchRate:  0.008766,
		},

		// Хобби и отдых  -> Спорт и отдых -> Фитнес и тренажёры
		{
			NodeID:               1054309,
			QueryFreq:            78475,
			QueryContactFreq:     1342,
			QuerySearchFreq:      38215,
			McatQueryFreq:        1775,
			McatQueryContactFreq: 48,
			McatQuerySearchFreq:  335,
			QueryMcatRate:        0.022619,
			QueryMcatContactRate: 0.035768,
			QueryMcatSearchRate:  0.008766,
		},

		// Хобби и отдых -> Спорт и отдых -> Фитнес и тренажёры -> Кардиотренажёры
		{
			NodeID:               2616267,
			QueryFreq:            78475,
			QueryContactFreq:     1342,
			QuerySearchFreq:      38215,
			McatQueryFreq:        1775,
			McatQueryContactFreq: 48,
			McatQuerySearchFreq:  335,
			QueryMcatRate:        0.022619,
			QueryMcatContactRate: 0.035768,
			QueryMcatSearchRate:  0.008766,
		},
		// Хобби и отдых  -> Спорт и отдых -> Бильярд и боулинг
		{
			NodeID:               1054306,
			QueryFreq:            78475,
			QueryContactFreq:     1342,
			QuerySearchFreq:      38215,
			McatQueryFreq:        1775,
			McatQueryContactFreq: 48,
			McatQuerySearchFreq:  335,
			QueryMcatRate:        0.022619,
			QueryMcatContactRate: 0.035768,
			QueryMcatSearchRate:  0.008766,
		},
	},
}

func TestQueryInfmNodesSubtreeBuilder_Build(t *testing.T) {
	cfg := envs.NewLocalTestConfig()
	path := cfg.Artifacts().InfmDictFallback()[infm.Default]
	normalizer := doppelgangerNormalizer.NewDoppelgangerNormalizer(map[string]string{})
	reader := dictionary.NewInfmDictReader(normalizer, infm.Default)

	_, tree, _ := reader.ReadInfmDict(dictionary.NewMockClient(path, infm.Default))
	localHierarchy := hierarchy.NewLocalHierarchy(tree)
	registry := hierarchy.NewMockRegistry([]int{1054309, 1054309, 2616267, 1054306})
	statsAccessor := catDict.NewNodeContactsAccessor()
	builder := NewQueryNodesSubtreeBuilder(localHierarchy, registry, statsAccessor)

	tests := []struct {
		name  string
		query string
		want  []int
	}{
		{
			name:  "Should Build correct node tree for 'бассейн'",
			query: "бассейн",
			want: []int{
				// root
				1047504,
				// Для дома и дачи -> Ремонт и строительство -> Сантехника, водоснабжение и сауна -> Товары для бани и сауны, бассейны -> Бассейны и комплектующие -> Бассейны
				1047737, 1047751, 1047754, 1288806, 2609066, 2618230,
				// Бассейны: Каркасная
				4698124,
				// Бассейны: Жёсткая
				4698122,
				// Бассейны: Bestway
				4698121,

				// Услуги -> Предложение услуг -> Красота, здоровье -> СПА-услуги, массаж
				1047505, 1047506, 1047554, 1047561,

				// Предложение услуг: Праздники, мероприятия -> Организация досуга и отдыха
				1047535, 4442348,

				// Недвижимость -> Дома, дачи, коттеджи -> Снять
				1054803, 1054804, 1054823,

				// Личные вещи -> Товары для детей и игрушки -> Товары для купания
				1054439, 1054789, 1054798,

				// Хобби и отдых -> Спорт и отдых
				1054252, 1054305,
			},
		},

		{
			name:  "Should Build correct node tree for 'диван бу'",
			query: "диван бу",
			want: []int{
				// root
				1047504,
				// Для дома и дачи -> Мебель и интерьер -> Кровати, диваны и кресла
				1047737, 1047738, 1047739,
				// Кровати, диваны и кресла: Диваны
				1090116,
				// Диваны: Прямые диваны -> Новые прямые диваны
				1288006, 2170275,
				// Диваны: Угловые диваны
				1288007,
				// Кровати, диваны и кресла: Кресла
				1210301,
			},
		},

		{
			name:  "Should Build correct node tree for 'query with terminator'",
			query: "query with terminator",
			want: []int{
				// root
				1047504,
				// Личные вещи
				1054439,
				// Хобби и отдых -> Спорт и отдых
				1054252, 1054305,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stats := queryStats[tt.query]
			dfsTraversalNodes := getDfsTraversalNodeIds(builder.Build(stats))

			if !reflect.DeepEqual(dfsTraversalNodes, tt.want) {
				t.Errorf("Build() = %v, want %v", dfsTraversalNodes, tt.want)
			}
		})
	}
}

func getDfsTraversalNodeIds(root *QueryStatsCategoryNode) []int {
	var result []int
	var dfs func(node *QueryStatsCategoryNode)
	dfs = func(node *QueryStatsCategoryNode) {
		if node == nil {
			return
		}

		result = append(result, node.Info.NodeID)
		for _, child := range node.Children {
			dfs(child)
		}
	}

	dfs(root)
	return result
}
