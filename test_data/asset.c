#include "ont.h"

char * init(){
    char * totalsupply = GetStorage("TCOIN_TOTAL_SUPPLY");
    if (arrayLen(totalsupply) > 0){
        return JsonMashalResult("this contract already initialized!","string");
    }else{
        PutStorage("TCOIN_TOTAL_SUPPLY","1000000000");
        PutStorage("00000001","1000000000");
    }
    return JsonMashalResult("init succeed","string");
}

char * getTotalSupply(){
    char * store = GetStorage("TCOIN_TOTAL_SUPPLY");
    if (arrayLen(store) > 0){
        int  totalsupply = Atoi(store);
        return JsonMashalResult(totalsupply,"int");
    }else{
        return JsonMashalResult("total supply has not been init!","string");
    }
}

char * balanceOf(char * address){
    int   balance =Atoi(GetStorage(address));
    return JsonMashalResult(balance,"int");
}

int transfer(char * from ,char * to, int amount){
    ContractLogError("------1-------");
    if(amount <= 0){
        return 0;
    }
    ContractLogError("------2-------");
    //checkwitness(from)
    if(strcmp(from,to) == 0){
       
        return 0;
    }
    ContractLogError("------3-------");
    char * fromValuestr = GetStorage(from);
        ContractLogError(from);
       ContractLogError("------4-------");
       ContractLogError(fromValuestr);
    int  fromValue = Atoi(fromValuestr);
       ContractLogError("------5-------");
    if (fromValue < amount){
        ContractLogError("------6-------");
        return 0;
    }
       ContractLogError("------7-------");
    if (fromValue == amount){
        DeleteStorage(from);
    }else{
           ContractLogError("------8-------");
        int  tovalue = Atoi(GetStorage(to));
           ContractLogError("------9-------");
        PutStorage(from,Itoa(fromValue -amount));
        PutStorage(to,Itoa(tovalue + amount));
    }
    return 1;
}

char * concat(char * a, char * b){
        int lena = arrayLen(a);
        int lenb = arrayLen(b);
        char * res = (char *)malloc((lena + lenb)*sizeof(char));
        for (int i = 0 ;i < lena ;i++){
                res[i] = a[i];
        }

        for (int j = 0; j < lenb ;j++){
                res[lenb + j] = b[j];
        }
        return res;
}

/*
*this is the common standard interface of ontology wasm contract
*/
char * invoke(char * method,char * args){

    if(strcmp(method,"init") == 0){
        char * result =  init();
        RuntimeNotify(result);
        return result;
    }

    if(strcmp(method,"totalSupply") == 0){

        char *result = getTotalSupply();
        RuntimeNotify(result);
        return result;
    }

    if(strcmp(method,"balanceOf") == 0){
        struct Param{
            char * address;
        };
        
        struct Param p;
        JsonUnmashalInput(&p,sizeof(p),args);

        char *result = balanceOf(p.address);
        RuntimeNotify(result);
        return result;
    }

    if(strcmp(method,"transfer") == 0){
        struct Param{
            char * from;
            char * to;
            int  amount;
        };
        
        struct Param p;
        JsonUnmashalInput(&p,sizeof(p),args);
        char * result;
        if(transfer(p.from,p.to,p.amount) > 0){

            result = JsonMashalResult(Itoa(p.amount),"string");
        }else{
            // result = JsonMashalResult(concat(concat(concat(concat(concat("transfer :",p.from),"to:"),p.to),I64toa(p.amount)),"failed!"),"string");
            result = JsonMashalResult(Itoa(p.amount),"string");
        }
        RuntimeNotify(result);
        return result;
    }
}