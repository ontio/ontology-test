void JsonUnmashalInput(void * addr,int size,char * arg);
char * JsonMashalResult(void * val,char * types);
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
        long long   totalsupply = Atoi64(store);
        return JsonMashalResult(totalsupply,"int64");
    }else{
        return JsonMashalResult("total supply has not been init!","string");
    }
}

char * balanceOf(char * address){
    long long  balance =Atoi64(GetStorage(address));
    return JsonMashalResult(balance,"int64");
}

int transfer(char * from ,char * to, int amount){
    if(amount <= 0){
        return 0;
    }
    //checkwitness(from)
    if(strcmp(from,to) == 0){

        return 0;
    }
    char * fromValuestr = GetStorage(from);
    long long  fromValue = Atoi64(fromValuestr);
    if (fromValue < amount){
        return 0;
    }
    if (fromValue == amount){
        DeleteStorage(from);
    }else{
        long long  tovalue = Atoi64(GetStorage(to));
        PutStorage(from,Itoa64(fromValue -amount));
        PutStorage(to,Itoa64(tovalue + amount));
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
        char * address = ReadStringParam(args);

        char *result = balanceOf(address);
        RuntimeNotify(result);
        return result;
    }

    if(strcmp(method,"transfer") == 0){

        char * from = ReadStringParam(args);
        char * to = ReadStringParam(args);
        long long amount = ReadInt64Param(args);

        char * result;
        if(transfer(from,to,amount) > 0){

            result = JsonMashalResult(Itoa64(amount),"string");
        }else{
            // result = JsonMashalResult(concat(concat(concat(concat(concat("transfer :",p.from),"to:"),p.to),I64toa(p.amount)),"failed!"),"string");
            result = JsonMashalResult(Itoa64(amount),"string");
        }
        RuntimeNotify(result);
        return result;
    }
}