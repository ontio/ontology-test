void JsonUnmashalInput(void * addr,int size,char * arg);
char * JsonMashalResult(void * val,char * types,int succeed);
int strcmp(char *a,char *b);
int arrayLen(char *a);
void * malloc(int size);
void RuntimeNotify(char * message);
void PutStorage(char * key,char *value);
char * GetStorage(char * key);
void DeleteStorage(char * key);
long long Atoi64(char * str);
char * I64toa(long long amount,int radix);
long Atoi(char * str);
char * Itoa(int a);
void ContractLogError(char *str);
int ReadInt32Param(char *addr);
long long ReadInt64Param(char *addr);
char * ReadStringParam(char *addr);
char * CallContract(char * address,char * contractCode,char * method,char * args);
char * MarshalNativeParams(void * args);
int CheckWitness(char * address);
char * GetSelfAddress();


char * init(){
    char * totalsupply = GetStorage("TCOIN_TOTAL_SUPPLY");
    if (arrayLen(totalsupply) > 0){
        return JsonMashalResult("this contract already initialized!","string",0);
    }else{
        PutStorage("TCOIN_TOTAL_SUPPLY","1000000000");
    }
    return JsonMashalResult("init succeed","string",1);
}

char * getTotalSupply(){
    char * store = GetStorage("TCOIN_TOTAL_SUPPLY");
    if (arrayLen(store) > 0){
        long long  totalsupply = Atoi64(store);
        return JsonMashalResult(totalsupply,"int64",1);
    }else{
        return JsonMashalResult("total supply has not been init!","string",0);
    }
}

char * balanceOf(char * address){
    long long   balance =Atoi64(GetStorage(address));
    return JsonMashalResult(balance,"int64",1);
}

int transfer(char * from ,char * to, long long amount){
    if(amount <= 0){
        return 0;
    }
    int witness = CheckWitness(from);
    if (witness == 0){
        return 0;
    }
    if(strcmp(from,to) == 0){

        return 0;
    }
    char * fromValuestr = GetStorage(from);

    long long fromValue = Atoi64(fromValuestr);
    if (fromValue < amount){
        return 0;
    }
    if (fromValue == amount){
        DeleteStorage(from);
    }else{
        long long  tovalue = Atoi64(GetStorage(to));
        PutStorage(from,I64toa(fromValue -amount,10));
        PutStorage(to,I64toa(tovalue + amount,10));
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

    char * args = MarshalNativeParams(transfer);
    char * result = CallContract("ff00000000000000000000000000000000000001","","transfer",args);
    if (strcmp(result,"true")==0){
        return 1;
    }else{
        return 0;
    }
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

        struct Param *p = (struct Param *)malloc(sizeof(struct Param));
        JsonUnmashalInput(p,sizeof(struct Param),args);

        char *result = balanceOf(p->address);
        RuntimeNotify(result);
        return result;
    }

    if(strcmp(method,"transfer") == 0){
        struct Param{
            char * from;
            char * to;
            long long  amount;
        };

        struct Param *p = (struct Param *)malloc(sizeof(struct Param));
        JsonUnmashalInput(p,sizeof(struct Param),args);
        char * result;
        if(transfer(p->from,p->to,p->amount) > 0){

            result = JsonMashalResult(Itoa(p->amount),"string",1);
        }else{
            result = JsonMashalResult(Itoa(p->amount),"string",0);
        }
        RuntimeNotify(result);
        return result;
    }
    if(strcmp(method,"collect") == 0){
        struct Param{
            char * from ;
            long long amount;
        };
        struct Param *p = (struct Param *)malloc(sizeof(struct Param));
        JsonUnmashalInput(p,sizeof(struct Param),args);

        int rate = 100;

        char * seflAddress = GetSelfAddress();
        ContractLogError(seflAddress);
        char * store = GetStorage("TCOIN_TOTAL_SUPPLY");
        long long  totalsupply = Atoi64(store);
        char * result = "";
        long long count = rate * p->amount;
        if (totalsupply >= count){
            int transOntResult = transOnt(p->from,seflAddress,p->amount);
            int tranCoin ;

            if (transOntResult == 1){
                PutStorage("TCOIN_TOTAL_SUPPLY",I64toa(totalsupply - count,10));

                char * toStore = GetStorage(p->from);
                long long  totalsupply = 0;
                if (arrayLen(toStore) > 0){
                    totalsupply = Atoi64(toStore);
                }

                PutStorage(p->from,I64toa(totalsupply + count,10));
                result = JsonMashalResult(I64toa(count,10),"int64",10);
            }else{
                result = JsonMashalResult("transfer ont failed","string",0);
            }
        }else{
            result = JsonMashalResult("not enough supply","string",0);
        }

        RuntimeNotify(result);
        return result;
    }
     if(strcmp(method,"withdraw") == 0){
        struct Param{
            long long amount;
        };
        struct Param *p = (struct Param *)malloc(sizeof(struct Param));
        JsonUnmashalInput(p,sizeof(struct Param),args);

        char * seflAddress = GetSelfAddress();
        ContractLogError(seflAddress);
        char * ownerAddress = "TA4ieHoEDmRmARQo6bVBayqPuvN51rd6wY";
        long long cnt = p->amount;
        int transOntResult = transOnt(seflAddress,ownerAddress,cnt);
        char * result = "";
        if (transOntResult ==1){
            result = JsonMashalResult(I64toa(cnt,10),"1nt64",1);
        }else{
            result = JsonMashalResult(I64toa(cnt,10),"1nt64",0);
        }

        RuntimeNotify(result);
        return result;
    }

}