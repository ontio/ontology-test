#include "ont.h"

char * invoke(char * method,char * args){

    if (strcmp(method,"getCurrentHeadHash")==0){
        char * hash = ONT_Block_GetCurrentHeaderHash();
        char * result = ONT_JsonMashalResult(hash,"string",1);
        ONT_Runtime_Notify(result);
        return result;
    }

    if(strcmp(method,"getCurrentHeaderHeight")==0){
        int height = ONT_Block_GetCurrentBlockHeight();
        char * result = ONT_JsonMashalResult(height,"int",1);
        ONT_Runtime_Notify(result);
        return result;
    }

    if(strcmp(method,"getCurrentBlockHash")==0){
        char * hash = ONT_Block_GetCurrentBlockHash();
        char * result = ONT_JsonMashalResult(hash,"string",1);
        ONT_Runtime_Notify(result);
        return result;
    }

    if(strcmp(method,"getCurrentBlockHeight")==0){
        int height = ONT_Block_GetCurrentBlockHeight();
        char * result = ONT_JsonMashalResult(height,"int",1);
        ONT_Runtime_Notify(result);
        return result;
    }

    if(strcmp(method,"getTransactionByHash")==0){
        char * hash = ONT_ReadStringParam(args);
        char * data = ONT_Block_GetTransactionByHash(hash);
        char * result = ONT_JsonMashalResult(data,"byte_array",1);
        ONT_Runtime_Notify(result);
        return result;
    }

    if(strcmp(method,"getTransactionCountByHash")==0){
        char * hash = ONT_ReadStringParam(args);
        int count = ONT_Block_GetTransactionCountByBlkHash(hash);
        char * result = ONT_JsonMashalResult(count,"int",1);
        ONT_Runtime_Notify(result);
        return result;
    }

    if(strcmp(method,"getTransactionCountByHeight")==0){
        int height = ONT_ReadInt32Param(args);
        int count = ONT_Block_GetTransactionCountByBlkHeight(height);
        char * result = ONT_JsonMashalResult(count,"int",1);
        ONT_Runtime_Notify(result);
        return result;
    }

    if(strcmp(method,"getTransactions")==0){
        char * hash = ONT_ReadStringParam(args);
        char ** trans = ONT_Block_GetTransactionsByBlkHash(hash);
        char * result = ONT_JsonMashalResult(trans,"string_array",1);
        ONT_Runtime_Notify(result);
        return result;
    }

    if(strcmp(method,"getTransactionsByHeight")==0){
        int height = ONT_ReadInt32Param(args);
        char ** trans = ONT_Block_GetTransactionsByBlkHeight(height);
        char * result = ONT_JsonMashalResult(trans,"string_array",1);
        ONT_Runtime_Notify(result);
        return result;
    }

    char * failed = ONT_JsonMashalResult(strconcat(method,"not supported"),"string",0);
    ONT_Runtime_Notify(failed);
    return failed;
}