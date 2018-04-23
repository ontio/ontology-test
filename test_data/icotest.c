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
        char * selfaddress = GetSelfAddress();
        PutStorage(selfaddress,"1000000000");
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
    int witness = Checkwitness(from);
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

int collect(char * from ,long long amount){
    char * selfaddress = GetSelfAddress();

    struct Param{
        char * from;
        char * to;
        long long amount;
    };
    struct Param *p =(struct Param *)malloc(sizeof(struct Param));
    p->from = from;
    p->to = selfaddress;
    p->amount = amount;

    char * result = CallContract("ff00000000000000000000000000000000000001","","transfer",p);
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
        JsonUnmashalInput(p,sizeof(p),args);

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
        JsonUnmashalInput(p,sizeof(p),args);
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
        JsonUnmashalInput(p,sizeof(p),args);

        int rate = 100;
        int transOnt = collect(p->from,p->amount);
        int tranCoin ;
        char * result;
        if (transOnt == 1){
            char * seflAddress = GetSelfAddress();
            tranCoin = transfer(seflAddress,p->from,p->amount * 100);
            if (tranCoin == 0){
                result = JsonMashalResult("transfer coin failed","string",0);
            }else{
                result = JsonMashalResult("transfer coin succeed","string",1);
            }
        }else{
            result = JsonMashalResult("transfer ont failed","string",0);
        }

        RuntimeNotify(result);
        return result;
    }

}