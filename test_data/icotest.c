#include "ont.h"

char * init(){
    char * totalsupply = ONT_Storage_Get("TCOIN_TOTAL_SUPPLY");
    if (arrayLen(totalsupply) > 0){
        return ONT_JsonMashalResult("this contract already initialized!","string",0);
    }else{
        ONT_Storage_Put("TCOIN_TOTAL_SUPPLY","1000000000");
    }
    return ONT_JsonMashalResult("init succeed","string",1);
}

char * getTotalSupply(){
    char * store = ONT_Storage_Get("TCOIN_TOTAL_SUPPLY");
    if (arrayLen(store) > 0){
        long long  totalsupply = Atoi64(store);
        return ONT_JsonMashalResult(totalsupply,"int64",1);
    }else{
        return ONT_JsonMashalResult("total supply has not been init!","string",0);
    }
}

char * balanceOf(char * address){
    long long   balance =Atoi64(ONT_Storage_Get(address));
    return ONT_JsonMashalResult(balance,"int64",1);
}

int transfer(char * from ,char * to, long long amount){
    if(amount <= 0){
        return 0;
    }
    int witness = ONT_Runtime_CheckWitness(from);
    if (witness == 0){
        return 0;
    }
    if(strcmp(from,to) == 0){

        return 0;
    }
    char * fromValuestr = ONT_Storage_Get(from);

    long long fromValue = Atoi64(fromValuestr);
    if (fromValue < amount){
        return 0;
    }
    if (fromValue == amount){
        ONT_Storage_Delete(from);
    }else{
        long long  tovalue = Atoi64(ONT_Storage_Get(to));
        ONT_Storage_Put(from,I64toa(fromValue -amount,10));
        ONT_Storage_Put(to,I64toa(tovalue + amount,10));
    }
    return 1;
}

int transOnt(char * from ,char * to,long long amount){

    struct State{
            int ver;
            char * from;
            char * to ;
            long long amount;
        };

    struct Transfer {
        int ver;
        struct State * states
    };


    struct State * state = (struct State*)malloc(sizeof(struct State));
    state->ver = 1;
    state->amount = amount;
    state->from = from;
    state->to = to;

    struct Transfer * transfer =(struct Transfer*)malloc(sizeof(struct Transfer));
    transfer->ver = 1;
    transfer->states = state;

    char * args = ONT_MarshalNativeParams(transfer);
    char * result = ONT_CallContract("ff00000000000000000000000000000000000001","","transfer",args);
    if (strcmp(result,"true")==0){
        return 1;
    }else{
        return 0;
    }
}

/*
*this is the common standard interface of ontology wasm contract
*/
char * invoke(char * method,char * args){

    if(strcmp(method,"init") == 0){
        char * result =  init();
        ONT_Runtime_Notify(result);
        return result;
    }

    if(strcmp(method,"totalSupply") == 0){

        char *result = getTotalSupply();
        ONT_Runtime_Notify(result);
        return result;
    }

    if(strcmp(method,"balanceOf") == 0){
        struct Param{
            char * address;
        };

        struct Param *p = (struct Param *)malloc(sizeof(struct Param));
        ONT_JsonUnmashalInput(p,sizeof(struct Param),args);

        char *result = balanceOf(p->address);
        ONT_Runtime_Notify(result);
        return result;
    }

    if(strcmp(method,"transfer") == 0){
        struct Param{
            char * from;
            char * to;
            long long  amount;
        };

        struct Param *p = (struct Param *)malloc(sizeof(struct Param));
        ONT_JsonUnmashalInput(p,sizeof(struct Param),args);
        char * result;
        if(transfer(p->from,p->to,p->amount) > 0){

            result = ONT_JsonMashalResult(Itoa(p->amount),"string",1);
        }else{
            result = ONT_JsonMashalResult(Itoa(p->amount),"string",0);
        }
        ONT_Runtime_Notify(result);
        return result;
    }
    if(strcmp(method,"collect") == 0){
        struct Param{
            char * from ;
            long long amount;
        };
        struct Param *p = (struct Param *)malloc(sizeof(struct Param));
        ONT_JsonUnmashalInput(p,sizeof(struct Param),args);

        int rate = 100;

        char * seflAddress = ONT_GetSelfAddress();
        char * store = ONT_Storage_Get("TCOIN_TOTAL_SUPPLY");
        long long  totalsupply = Atoi64(store);
        char * result = "";
        long long count = rate * p->amount;
        if (totalsupply >= count){
            int transOntResult = transOnt(p->from,seflAddress,p->amount);
            int tranCoin ;

            if (transOntResult == 1){
                ONT_Storage_Put("TCOIN_TOTAL_SUPPLY",I64toa(totalsupply - count,10));

                char * toStore = ONT_Storage_Get(p->from);
                long long  tovalue = 0;
                if (arrayLen(toStore) > 0){
                    tovalue = Atoi64(toStore);
                }

                ONT_Storage_Put(p->from,I64toa(tovalue + count,10));
                result = ONT_JsonMashalResult(count,"int64",1);
            }else{
                result = ONT_JsonMashalResult("transfer ont failed","string",0);
            }
        }else{
            result = ONT_JsonMashalResult("not enough supply","string",0);
        }

        ONT_Runtime_Notify(result);
        return result;
    }
     if(strcmp(method,"withdraw") == 0){
        struct Param{
            long long amount;
        };
        struct Param *p = (struct Param *)malloc(sizeof(struct Param));
        ONT_JsonUnmashalInput(p,sizeof(struct Param),args);

        char * seflAddress = ONT_GetSelfAddress();
        char * ownerAddress = "TA4ieHoEDmRmARQo6bVBayqPuvN51rd6wY";
        long long cnt = p->amount;
        int transOntResult = transOnt(seflAddress,ownerAddress,cnt);
        char * result = "";
        if (transOntResult ==1){
            result = ONT_JsonMashalResult(cnt,"int64",1);
        }else{
            result = ONT_JsonMashalResult(cnt,"int64",0);
        }

        ONT_Runtime_Notify(result);
        return result;
    }

}