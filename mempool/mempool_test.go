package mempool

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestMemPool_ReadAndWrite(t *testing.T) {
	tests := map[string]struct {
		input string
		expectedOutput string
		expectedError string
	} {
		"valid input": {
			input: `TxHash=40E10C7CF56A738C0B8AD4EE30EA8008C7B2334B3ADA195083F8CB18BD3911A0 Gas=729000 FeePerGas=0.11134106816568039 Signature=6386A3893BEB6A5A64E0677F406634E791DEE78D49CF30581AE5281D4094E495E671647EF5E7FD2D207AB8EBA0EA693703E9C368402731BE99E81BDB748EA662
TxHash=4B2B252899DC689106C8FCEA3E24E4AFFC597D2B4E701F99EB8CD909217D323F Gas=834000 FeePerGas=0.27503931836911927 Signature=88B520FC81B8F8D1FD7B1F42B38481426CB5CE7C27F4A03F51C4A4710A0DC5FA3127E3DE20A818555CE74470A3420E39F65FE7D5053FBBC7C2151A3F22081A5B
TxHash=F75F133F149FDA7DEB391B2446C5196E7C704F45456E69312C310C72893F5B6A Gas=303000 FeePerGas=0.13954792852514392 Signature=E7CF053A8E21351EFE0C52869AFDDD79BF57E92AE918C38ECA8744FD640B1D8D7FF40DF296C2280C71FA78BF5F2A72E33D17C279D30CEAFC7F205521D2B96678
TxHash=16633D0A25ECA886F100A34BA5C43366732836E6E7B140159298C71CF78309F9 Gas=179000 FeePerGas=0.6767094874229113 Signature=0B27A09EFFD07C035152EE4431066F2418735F4BAEB290E9F099EA0967133F471C10C21E86E8F28B674E2CB67E4A489BB46B7BCA44D8C670928CCF8A3567D443`,
			expectedOutput: `TxHash=4B2B252899DC689106C8FCEA3E24E4AFFC597D2B4E701F99EB8CD909217D323F Gas=834000 FeePerGas=0.27503931836911927 Signature=88B520FC81B8F8D1FD7B1F42B38481426CB5CE7C27F4A03F51C4A4710A0DC5FA3127E3DE20A818555CE74470A3420E39F65FE7D5053FBBC7C2151A3F22081A5B
TxHash=16633D0A25ECA886F100A34BA5C43366732836E6E7B140159298C71CF78309F9 Gas=179000 FeePerGas=0.6767094874229113 Signature=0B27A09EFFD07C035152EE4431066F2418735F4BAEB290E9F099EA0967133F471C10C21E86E8F28B674E2CB67E4A489BB46B7BCA44D8C670928CCF8A3567D443
TxHash=40E10C7CF56A738C0B8AD4EE30EA8008C7B2334B3ADA195083F8CB18BD3911A0 Gas=729000 FeePerGas=0.11134106816568039 Signature=6386A3893BEB6A5A64E0677F406634E791DEE78D49CF30581AE5281D4094E495E671647EF5E7FD2D207AB8EBA0EA693703E9C368402731BE99E81BDB748EA662
TxHash=F75F133F149FDA7DEB391B2446C5196E7C704F45456E69312C310C72893F5B6A Gas=303000 FeePerGas=0.13954792852514392 Signature=E7CF053A8E21351EFE0C52869AFDDD79BF57E92AE918C38ECA8744FD640B1D8D7FF40DF296C2280C71FA78BF5F2A72E33D17C279D30CEAFC7F205521D2B96678`,
		},
		"valid input with whitespaces": {
			input: `
TxHash=40E10C7CF56A738C0B8AD4EE30EA8008C7B2334B3ADA195083F8CB18BD3911A0 Gas=729000 FeePerGas=0.11134106816568039 Signature=6386A3893BEB6A5A64E0677F406634E791DEE78D49CF30581AE5281D4094E495E671647EF5E7FD2D207AB8EBA0EA693703E9C368402731BE99E81BDB748EA662
  		
TxHash=4B2B252899DC689106C8FCEA3E24E4AFFC597D2B4E701F99EB8CD909217D323F Gas=834000 FeePerGas=0.27503931836911927 Signature=88B520FC81B8F8D1FD7B1F42B38481426CB5CE7C27F4A03F51C4A4710A0DC5FA3127E3DE20A818555CE74470A3420E39F65FE7D5053FBBC7C2151A3F22081A5B
TxHash=F75F133F149FDA7DEB391B2446C5196E7C704F45456E69312C310C72893F5B6A Gas=303000 FeePerGas=0.13954792852514392 Signature=E7CF053A8E21351EFE0C52869AFDDD79BF57E92AE918C38ECA8744FD640B1D8D7FF40DF296C2280C71FA78BF5F2A72E33D17C279D30CEAFC7F205521D2B96678
TxHash=16633D0A25ECA886F100A34BA5C43366732836E6E7B140159298C71CF78309F9 Gas=179000 FeePerGas=0.6767094874229113 Signature=0B27A09EFFD07C035152EE4431066F2418735F4BAEB290E9F099EA0967133F471C10C21E86E8F28B674E2CB67E4A489BB46B7BCA44D8C670928CCF8A3567D443
`,
			expectedOutput: `TxHash=4B2B252899DC689106C8FCEA3E24E4AFFC597D2B4E701F99EB8CD909217D323F Gas=834000 FeePerGas=0.27503931836911927 Signature=88B520FC81B8F8D1FD7B1F42B38481426CB5CE7C27F4A03F51C4A4710A0DC5FA3127E3DE20A818555CE74470A3420E39F65FE7D5053FBBC7C2151A3F22081A5B
TxHash=16633D0A25ECA886F100A34BA5C43366732836E6E7B140159298C71CF78309F9 Gas=179000 FeePerGas=0.6767094874229113 Signature=0B27A09EFFD07C035152EE4431066F2418735F4BAEB290E9F099EA0967133F471C10C21E86E8F28B674E2CB67E4A489BB46B7BCA44D8C670928CCF8A3567D443
TxHash=40E10C7CF56A738C0B8AD4EE30EA8008C7B2334B3ADA195083F8CB18BD3911A0 Gas=729000 FeePerGas=0.11134106816568039 Signature=6386A3893BEB6A5A64E0677F406634E791DEE78D49CF30581AE5281D4094E495E671647EF5E7FD2D207AB8EBA0EA693703E9C368402731BE99E81BDB748EA662
TxHash=F75F133F149FDA7DEB391B2446C5196E7C704F45456E69312C310C72893F5B6A Gas=303000 FeePerGas=0.13954792852514392 Signature=E7CF053A8E21351EFE0C52869AFDDD79BF57E92AE918C38ECA8744FD640B1D8D7FF40DF296C2280C71FA78BF5F2A72E33D17C279D30CEAFC7F205521D2B96678`,
		},
		"empty input": {
			input: `
	   
  `,
			expectedOutput: "",
		},
		"invalid input": {
			input: "abcd",
			expectedError: "Invalid token [abcd]",
		},
	}

	for tName, tc := range tests {
		tc := tc
		t.Run(tName, func(t *testing.T) {
			memPool := NewMemPool()
			reader := strings.NewReader(tc.input)
			err := memPool.ReadTransactions(reader)
			if tc.expectedError != "" {
				require.EqualError(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)
				var buf bytes.Buffer
				err = memPool.WriteTransactions(&buf)
				require.NoError(t, err)
				output := buf.String()
				require.Equal(t, tc.expectedOutput, output)
			}
		})
	}
}