//system apis
void * calloc(int count,int length);
void * malloc(int size);
int arrayLen(void *a);
int memcpy(void * dest,void * src,int length);
int memset(void * dest,char c,int length);

//utility apis
int strcmp(char *a,char *b);
void * strconcat(char *array);
int Atoi(char * s);
long long Atoi64(char *s);
char * Itoa(int a);
char * I64toa(long long amount,int radix);
char * SHA1(char *s);
char * SHA256(char *s);

//parameter apis
int ONT_ReadInt32Param(char *args);
long long ONT_ReadInt64Param(char * args);
char * ONT_ReadStringParam(char * args);
void ONT_JsonUnmashalInput(void * addr,int size,char * arg);
char * ONT_JsonMashalResult(void * val,char * types,int succeed);
char * ONT_JsonMashalParams(void * s);
char * ONT_RawMashalParams(void *s);
char * ONT_GetCallerAddress();
char * ONT_GetSelfAddress();
char * ONT_CallContract(char * address,char * contractCode,char * method,char * args);
char * ONT_MarshalNativeParams(void *);

//Runtime apis
int ONT_Runtime_CheckWitness(char * address);
void ONT_Runtime_Notify(char * address);
int ONT_Runtime_CheckSig(char * pubkey,char * data,char * sig);
int ONT_Runtime_GetTime();
void ONT_Runtime_Log(char * message);

//Attribute apis
int ONT_Attribute_GetUsage(char * data);
char * ONT_Attribute_GetData(char * data);

//Block apis
char * ONT_Block_GetCurrentHeaderHash();
int ONT_Block_GetCurrentHeaderHeight();
char * ONT_Block_GetCurrentBlockHash();
int ONT_Block_GetCurrentBlockHeight();
char * ONT_Block_GetTransactionByHash(char * hash);
char * ONT_Block_GetTransactionCount(int count);
char * ONT_Block_GetTransactions(char * data);

//Blockchain apis
int ONT_BlockChain_GetHeight();
char * ONT_BlockChain_GetHeaderByHeight(int height);
char * ONT_BlockChain_GetHeaderByHash(char * hash);
char * ONT_BlockChain_GetBlockByHeight(int height);
char * ONT_BlockChain_GetBlockByHash(char * hash);
char * ONT_BlockChain_GetContract(char * address);

//header apis
char * ONT_Header_GetHash(char * data);
int ONT_Header_GetVersion(char * data);
char * ONT_Header_GetPrevHash(char * data);
char * ONT_Header_GetMerkleRoot(char  * data);
int ONT_Header_GetIndex(char * data);
int ONT_Header_GetTimestamp(char * data);
long long ONT_Header_GetConsensusData(char * data);
char * ONT_Header_GetNextConsensus(char * data);

//storage apis
void ONT_Storage_Put(char * key,char * value);
char * ONT_Storage_Get(char * key);
void ONT_Storage_Delete(char * key);

//transaction apis
char * ONT_Transaction_GetHash(char * data);
int ONT_Transaction_GetType(char * data);
char * ONT_Transaction_GetAttributes(char * data);
/*
*this is the common standard interface of ontology wasm contract
*/
char * invoke(char * method,char * args){

    if(strcmp(method,"transferont") == 0){
        char * from = ONT_ReadStringParam(args);
        char * to = ONT_ReadStringParam(args);
        long long ontAmount = ONT_ReadInt64Param(args);
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
        state->amount = ontAmount;
        state->from = from;
        state->to = to;

        struct Transfer * transfer =(struct Transfer*)malloc(sizeof(struct Transfer));
        transfer->ver = 1;
        transfer->states = state;
        //1. call native contract to tranfer ont

        char * args = ONT_MarshalNativeParams(transfer);
        char * result = ONT_CallContract("ff00000000000000000000000000000000000001","","transfer",args);
        char * res;
        if (strcmp(result,"true")==0){
           res = ONT_JsonMashalResult("transfer succeed","string",1);
        }else{
           res = ONT_JsonMashalResult("transfer failed","string",0);
        }
        ONT_Runtime_Notify(res);
        return res;
    }
}