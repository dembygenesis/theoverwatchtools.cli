package urlutil

//func TestMapToQueryString(t *testing.T) {
//	tests := []struct {
//		name    string
//		baseUrl string
//		params  map[string][]string
//		want    string
//		wantErr bool
//	}{
//		{
//			name:    "Append to URL without query",
//			baseUrl: "http://example.com/path",
//			params: map[string][]string{
//				"key1": {"value1", "value2"},
//				"key2": {"value3"},
//			},
//			want:    "http://example.com/path?key1=value1&key1=value2&key2=value3",
//			wantErr: false,
//		},
//		{
//			name:    "Append to a partial URL without query",
//			baseUrl: "/v2/path",
//			params: map[string][]string{
//				"key1": {"value1", "value2"},
//				"key2": {"value3"},
//			},
//			want:    "/v2/path?key1=value1&key1=value2&key2=value3",
//			wantErr: false,
//		},
//		{
//			name:    "Append to URL with existing query",
//			baseUrl: "http://example.com/path?existing=param",
//			params: map[string][]string{
//				"key1": {"value1", "value2"},
//			},
//			want:    "http://example.com/path?existing=param&key1=value1&key1=value2",
//			wantErr: false,
//		},
//		{
//			name:    "Special characters handling",
//			baseUrl: "http://example.com/path",
//			params: map[string][]string{
//				"key 1": {"value 1", "value&2"},
//			},
//			want:    "http://example.com/path?key+1=value+1&key+1=value%262",
//			wantErr: false,
//		},
//		{
//			name:    "Invalid base URL",
//			baseUrl: "http://a b.com/",
//			params: map[string][]string{
//				"key": {"value"},
//			},
//			wantErr: true,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := MapToQueryString(tt.baseUrl, tt.params)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("MapToQueryString() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if got != tt.want {
//				t.Errorf("MapToQueryString() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
