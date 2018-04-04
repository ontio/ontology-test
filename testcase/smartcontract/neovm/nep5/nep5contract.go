/*
 * Copyright (C) 2018 The ontology Authors
 * This file is part of The ontology library.
 *
 * The ontology is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The ontology is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The ontology.  If not, see <http://www.gnu.org/licenses/>.
 */
package nep5

import "github.com/ontio/ontology-go-sdk/utils"

// NEP-5 smart contract hex string
var nep5Contract = "0111c56b6c766b00527ac46c766b51527ac4616168164e656f2e52756e74696d652e47657454726967676572609c6c766b5" +
	"2527ac46c766b52c364cb01616c766b00c304696e6974876c766b53527ac46c766b53c36411006165e4016c766b54527ac462ad016c766b00c" +
	"30b746f74616c537570706c79876c766b55527ac46c766b55c364110061658e036c766b54527ac4627e016c766b00c3046e616d65876c766b5" +
	"6527ac46c766b56c3641100616569016c766b54527ac46256016c766b00c30673796d626f6c876c766b57527ac46c766b57c36411006165500" +
	"16c766b54527ac4622c016c766b00c3087472616e73666572876c766b58527ac46c766b58c3647100616c766b51c3c0539c009c6c766b5c527" +
	"ac46c766b5cc3640e00006c766b54527ac462e9006c766b51c300c36c766b59527ac46c766b51c351c36c766b5a527ac46c766b51c352c36c7" +
	"66b5b527ac46c766b59c36c766b5ac36c766b5bc361527265ff026c766b54527ac462a0006c766b00c30962616c616e63654f66876c766b5d5" +
	"27ac46c766b5dc3644900616c766b51c3c0519c009c6c766b5f527ac46c766b5fc3640e00006c766b54527ac4625c006c766b51c300c36c766" +
	"b5e527ac46c766b5ec36165df046c766b54527ac4623b006c766b00c308646563696d616c73876c766b60527ac46c766b60c364110061653e0" +
	"06c766b54527ac4620f0061006c766b54527ac46203006c766b54c3616c756600c56b094e594320546f6b656e616c756600c56b034e5943616" +
	"c756600c56b58616c756653c56b616168164e656f2e53746f726167652e476574436f6e746578740b746f74616c537570706c79617c680f4e6" +
	"56f2e53746f726167652e4765746c766b00527ac46c766b00c3c000a06c766b51527ac46c766b51c3640e00006c766b52527ac4626e0161681" +
	"64e656f2e53746f726167652e476574436f6e74657874611401e2d4a41080e8097758c53f74b3081ea2778b540700003426f56b1c615272680" +
	"f4e656f2e53746f726167652e507574616100611401e2d4a41080e8097758c53f74b3081ea2778b540700003426f56b1c615272087472616e7" +
	"366657254c168124e656f2e52756e74696d652e4e6f74696679616168164e656f2e53746f726167652e476574436f6e74657874611401f1200" +
	"8b5700537aea16594c26861d15cf4e5ae070080f420e6b500615272680f4e656f2e53746f726167652e507574616100611401f12008b570053" +
	"7aea16594c26861d15cf4e5ae070080f420e6b500615272087472616e7366657254c168124e656f2e52756e74696d652e4e6f7469667961616" +
	"8164e656f2e53746f726167652e476574436f6e746578740b746f74616c537570706c79070000c16ff28623615272680f4e656f2e53746f726" +
	"167652e50757461516c766b52527ac46203006c766b52c3616c756651c56b616168164e656f2e53746f726167652e476574436f6e746578740" +
	"b746f74616c537570706c79617c680f4e656f2e53746f726167652e4765746c766b00527ac46203006c766b00c3616c75665bc56b6c766b005" +
	"27ac46c766b51527ac46c766b52527ac4616c766b52c300a16c766b55527ac46c766b55c3640e00006c766b56527ac46205026c766b00c3616" +
	"8184e656f2e52756e74696d652e436865636b5769746e657373009c6c766b57527ac46c766b57c3640e00006c766b56527ac462c9016c766b0" +
	"0c36c766b51c39c6c766b58527ac46c766b58c3640e00516c766b56527ac462a4016168164e656f2e53746f726167652e476574436f6e74657" +
	"8746c766b00c3617c680f4e656f2e53746f726167652e4765746c766b53527ac46c766b53c36c766b52c39f6c766b59527ac46c766b59c3640" +
	"e00006c766b56527ac46247016c766b53c36c766b52c39c6c766b5a527ac46c766b5ac3643b006168164e656f2e53746f726167652e4765744" +
	"36f6e746578746c766b00c3617c68124e656f2e53746f726167652e44656c657465616241006168164e656f2e53746f726167652e476574436" +
	"f6e746578746c766b00c36c766b53c36c766b52c394615272680f4e656f2e53746f726167652e507574616168164e656f2e53746f726167652" +
	"e476574436f6e746578746c766b51c3617c680f4e656f2e53746f726167652e4765746c766b54527ac46168164e656f2e53746f726167652e4" +
	"76574436f6e746578746c766b51c36c766b54c36c766b52c393615272680f4e656f2e53746f726167652e50757461616c766b00c36c766b51c" +
	"36c766b52c3615272087472616e7366657254c168124e656f2e52756e74696d652e4e6f7469667961516c766b56527ac46203006c766b56c36" +
	"16c756652c56b6c766b00527ac4616168164e656f2e53746f726167652e476574436f6e746578746c766b00c3617c680f4e656f2e53746f726" +
	"167652e4765746c766b51527ac46203006c766b51c3616c7566"

	//Nep5 contract address
var nep5Address = utils.GetNeoVMContractAddress(nep5Contract)

/*
using Neo.SmartContract.Framework;
using Neo.SmartContract.Framework.Services.Neo;
using Neo.SmartContract.Framework.Services.System;
using System;
using System.ComponentModel;
using System.Numerics;

namespace Ontology
{
    public class Ontology : SmartContract
    {
        //Token Settings
        public static string Name() => "NYC Token";
        public static string Symbol() => "NYC";
        public static readonly byte[] crowdSale = { 1, 241, 32, 8, 181, 112, 5, 55, 174, 161, 101, 148, 194, 104, 97, 209, 92, 244, 229, 174 };
        public static readonly byte[] institution = {1, 226 ,212, 164, 16, 128, 232, 9, 119, 88, 197, 63, 116, 179, 8, 30, 162, 119, 139, 84};
        //public static readonly byte[] community = "TA8EzSQpa8WmjGwD2cxcYAEQRErUqyVaZA".ToScriptHash();
        public static byte Decimals() => 8;
        private const ulong factor = 100000000; //decided by Decimals()

        private const ulong totalAmount = 100000000 * factor;
        private const ulong crowdSaleCap = 2000000 * factor;
        private const ulong institutionCap =80000000 * factor;
        //private const ulong communityCap = 50000000 * factor;

        [DisplayName("transfer")]
        public static event Action<byte[], byte[], BigInteger> Transferred;

        public static Object Main(string operation, params object[] args)
        {
            if (Runtime.Trigger == TriggerType.Application)
            {
                if (operation == "init") return Init();
                if (operation == "totalSupply") return TotalSupply();
                if (operation == "name") return Name();
                if (operation == "symbol") return Symbol();
                if (operation == "transfer")
                {
                    if (args.Length != 3) return false;
                    byte[] from = (byte[])args[0];
                    byte[] to = (byte[])args[1];
                    BigInteger value = (BigInteger)args[2];
                    return Transfer(from, to, value);
                }
                if (operation == "balanceOf")
                {
                    if (args.Length != 1) return 0;
                    byte[] account = (byte[])args[0];
                    return BalanceOf(account);
                }
                if (operation == "decimals") return Decimals();
            }
            return false;
        }

        public static bool Init()
        {
            byte[] total_supply = Storage.Get(Storage.CurrentContext, "totalSupply");
            if (total_supply.Length != 0) return false;

            //Storage.Put(Storage.CurrentContext, community, communityCap);
            //Transferred(null, community, communityCap);

            Storage.Put(Storage.CurrentContext, institution, institutionCap);
            Transferred(null, institution, institutionCap);

            Storage.Put(Storage.CurrentContext, crowdSale, crowdSaleCap);
            Transferred(null, crowdSale, crowdSaleCap);

            Storage.Put(Storage.CurrentContext, "totalSupply", totalAmount);
            return true;
        }

        // get the total supplied token
        public static BigInteger TotalSupply()
        {
            return Storage.Get(Storage.CurrentContext, "totalSupply").AsBigInteger();
        }

        // function that is always called when someone wants to transfer tokens.
        public static bool Transfer(byte[] from, byte[] to, BigInteger value)
        {
            if (value <= 0) return false;
            if (!Runtime.CheckWitness(from)) return false;
            if (from == to) return true;
            BigInteger from_value = Storage.Get(Storage.CurrentContext, from).AsBigInteger();
            if (from_value < value) return false;
            if (from_value == value)
                Storage.Delete(Storage.CurrentContext, from);
            else
                Storage.Put(Storage.CurrentContext, from, from_value - value);
            BigInteger to_value = Storage.Get(Storage.CurrentContext, to).AsBigInteger();
            Storage.Put(Storage.CurrentContext, to, to_value + value);
            Transferred(from, to, value);
            return true;
        }

        // get token amount of specific address
        public static BigInteger BalanceOf(byte[] address)
        {
            return Storage.Get(Storage.CurrentContext, address).AsBigInteger();
        }
    }
}
*/
