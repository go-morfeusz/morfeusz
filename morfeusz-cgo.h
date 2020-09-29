#ifndef MORFEUSZ_CGO_H
#define MORFEUSZ_CGO_H

#ifdef __cplusplus
extern "C" {
#endif  // __cplusplus

typedef void* Morf;
typedef void* Res;
// C++ code has to communicate with Go code via C.
// Struct String is an intermediary type between
// std::string and Go string and vice versa.
struct String {
    const char* p;
    int n;
};
typedef struct String Error;
struct StringArray {
    const struct String* strings;
    int length;
};
struct TokenInfo {
    struct String orth;
    struct String lemma;
    int startNode;
    int endNode;
    int tagID;
    int nameID;
    int labelsID;
};
struct TokenInfoArray {
    const struct TokenInfo* tokens;
    int length;
    Error error;
};
enum Charset {
    UTF8,
    ISO8859_2,
    CP1250,
    CP852
};
enum TokenNumbering {
    SEPARATE_NUMBERING,
    CONTINUOUS_NUMBERING
};
enum CaseHandling {
    CONDITIONALLY_CASE_SENSITIVE,
    STRICTLY_CASE_SENSITIVE,
    IGNORE_CASE
};
enum WhitespaceHandling {
    SKIP_WHITESPACES,
    APPEND_WHITESPACES,
    KEEP_WHITESPACES
};
enum Usage {
    BOTH_ANALYSE_AND_GENERATE,
    ANALYSE_ONLY,
    GENERATE_ONLY
};
Morf createInstance(const struct String dictName, enum Usage usage);
Res analyseString(const Morf m, const struct String text);
int hasNext(Res r);
const struct TokenInfo next(Res r);
const struct String tagsetId(const Morf m);
const struct String tag(const Morf m, int tagId);
int tagId(const Morf m, const struct String tag);
const struct String name(const Morf m, int nameId);
int nameId(const Morf m, const struct String name);
const struct String labelsAsString(const Morf m, int labelsId);
const struct StringArray labels(const Morf m, int labelsId);
int labelsId(const Morf m, const struct String labels);
int tagsCount(const Morf m);
int namesCount(const Morf m);
int labelsCount(const Morf m);
const struct TokenInfoArray generate(const Morf m, const struct String lemma);
const struct TokenInfoArray generateWithTagID(
    const Morf m, int tagId, const struct String lemma);
const struct String dictId(const Morf m);
const struct String dictCopyright(const Morf m);
const Error setAggl(Morf m, const struct String aggl);
const Error setPraet(Morf m, const struct String praet);
const Error setCharset(Morf m, enum Charset encoding);
const Error setCaseHandling(Morf m, enum CaseHandling caseHandling);
const Error setTokenNumbering(Morf m, enum TokenNumbering numbering);
const Error setWhitespaceHandling(Morf m, enum WhitespaceHandling handling);
const Error setDictionary(Morf m, struct String dictName);
void setDebug(Morf m, int debug);
const struct String aggl(const Morf m);
const struct String praet(const Morf m);
enum Charset charset(const Morf m);
enum CaseHandling caseHandling(const Morf m);
enum TokenNumbering tokenNumbering(const Morf m);
enum WhitespaceHandling whitespaceHandling(const Morf m);
const struct StringArray availableAgglOptions(const Morf m);
const struct StringArray availablePraetOptions(const Morf m);
const struct StringArray dictionarySearchPaths(const Morf m);
void prependToDictionarySearchPaths(Morf m, const struct String path);
void appendToDictionarySearchPaths(Morf m, const struct String path);
int removeFromDictionarySearchPaths(Morf m, const struct String path);
void clearDictionarySearchPaths(Morf m);
Morf cloneMorf(const Morf m);
const struct String version(void);
const struct String defaultDictName(void);
const struct String copyright(void);
void freeMorf(const Morf m);
void freeRes(const Res r);
void freeTokenInfo(const struct TokenInfo* t);
void freeStringArray(const struct StringArray* arr);
void freeTokenInfoArray(const struct TokenInfoArray* arr);
void freeCharArray(const char* p);

#ifdef __cplusplus
}  // extern "C"
#endif  // __cplusplus

#endif // MORFEUSZ_CGO_H
