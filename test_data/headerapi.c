#include "ont.h"

char * invoke(char * method,char * args){

    if(strcmp(method,"getHeaderHashByHeight")==0){
        int height = ONT_ReadInt32Param(args);
        char * hash = ONT_Header_GetHashByHeight(height);
        char * result = ONT_JsonMashalResult(hash,"string",1);
        ONT_Runtime_Notify(result);
        return result;
    }

    if(strcmp(method,"getHeaderVersionByHeight")==0){
        int height = ONT_ReadInt32Param(args);
        int version = ONT_Header_GetVersionByHeight(height);
        char * result = ONT_JsonMashalResult(version,"int",1);
        ONT_Runtime_Notify(result);
        return result;
    }

    if(strcmp(method,"getHeaderVersionByHash")==0){
        char * hash = ONT_ReadStringParam(args);
        int version = ONT_Header_GetVersionByHash(hash);
        char * result = ONT_JsonMashalResult(version,"int",1);
        ONT_Runtime_Notify(result);
        return result;
    }

    if(strcmp(method,"getPrevHashByHeight")==0){
        int height = ONT_ReadInt32Param(args);
        char * prevHash = ONT_Header_GetPrevHashByHeight(height);
        char * result = ONT_JsonMashalResult(prevHash,"string",1);
        ONT_Runtime_Notify(result);
        return result;
    }

    if(strcmp(method,"getPrevHashByHash")==0){
        char * hash = ONT_ReadStringParam(args);
        char * prevHash = ONT_Header_GetPrevHashByHash(hash);
        char * result = ONT_JsonMashalResult(prevHash,"string",1);
        ONT_Runtime_Notify(result);
        return result;
    }

    if(strcmp(method,"getMerkelRootByHeight")==0){
        int height = ONT_ReadInt32Param(args);
        char * prevHash = ONT_Header_GetMerkleRootByHeight(height);
        char * result = ONT_JsonMashalResult(prevHash,"string",1);
        ONT_Runtime_Notify(result);
        return result;
    }

    if(strcmp(method,"getMerkelRootByHash")==0){
        char * hash = ONT_ReadStringParam(args);
        char * prevHash = ONT_Header_GetMerkleRootByHash(hash);
        char * result = ONT_JsonMashalResult(prevHash,"string",1);
        ONT_Runtime_Notify(result);
        return result;
    }

    if(strcmp(method,"getTimestampByHeight")==0){
        int height = ONT_ReadInt32Param(args);
        int timestamp = ONT_Header_GetTimestampByHeight(height);
        char * result = ONT_JsonMashalResult(timestamp,"int",1);
        ONT_Runtime_Notify(result);
        return result;
    }

    if(strcmp(method,"getTimestampByHash")==0){
        char * hash = ONT_ReadStringParam(args);
        int timestamp = ONT_Header_GetTimestampByHash(hash);
        char * result = ONT_JsonMashalResult(timestamp,"int",1);
        ONT_Runtime_Notify(result);
        return result;
    }

    if(strcmp(method,"getIndexByHash")==0){
        char * hash = ONT_ReadStringParam(args);
        int index = ONT_Header_GetIndexByHash(hash);
        char * result = ONT_JsonMashalResult(index,"int",1);
        ONT_Runtime_Notify(result);
        return result;
    }

    if(strcmp(method,"getConsensusDataByHeight")==0){
        int height = ONT_ReadInt32Param(args);
        char * data = ONT_Header_GetConsensusDataByHeight(height);
        char * result = ONT_JsonMashalResult(data,"int64",1);
        ONT_Runtime_Notify(result);
        return result;
    }

    if(strcmp(method,"getConsensusDataByHash")==0){
        char * hash = ONT_ReadStringParam(args);
        char * data = ONT_Header_GetConsensusDataByHash(hash);
        char * result = ONT_JsonMashalResult(data,"int64",1);
        ONT_Runtime_Notify(result);
        return result;
    }

    if(strcmp(method,"getNextConsensusByHeight")==0){
        int height = ONT_ReadInt32Param(args);
        char * data = ONT_Header_GetNextConsensusByHeight(height);
        char * result = ONT_JsonMashalResult(data,"string",1);
        ONT_Runtime_Notify(result);
        return result;
    }

    if(strcmp(method,"getNextConsensusByHash")==0){
        char * hash = ONT_ReadStringParam(args);
        char * data = ONT_Header_GetNextConsensusByHash(hash);
        char * result = ONT_JsonMashalResult(data,"string",1);
        ONT_Runtime_Notify(result);
        return result;
    }


    char * failed = ONT_JsonMashalResult(strconcat(method,"not supported"),"string",0);
    ONT_Runtime_Notify(failed);
    return failed;
}
