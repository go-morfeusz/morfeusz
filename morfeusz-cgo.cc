#include "morfeusz-cgo.h"
#include "morfeusz2.h"

#include <string.h>
#include <exception>
#include <list>
#include <set>
#include <string>
#include <vector>

using morfeusz::Morfeusz;
using morfeusz::ResultsIterator;
using morfeusz::MorphInterpretation;
using morfeusz::IdResolver;

namespace {

const morfeusz::Charset translateCharset[] = {
    morfeusz::Charset::UTF8,
    morfeusz::Charset::ISO8859_2,
    morfeusz::Charset::CP1250,
    morfeusz::Charset::CP852,
};
const morfeusz::TokenNumbering translateTokenNumbering[] = {
    morfeusz::TokenNumbering::SEPARATE_NUMBERING,
    morfeusz::TokenNumbering::CONTINUOUS_NUMBERING,
};
const morfeusz::CaseHandling translateCaseHandling[] = {
    morfeusz::CaseHandling::CONDITIONALLY_CASE_SENSITIVE,
    morfeusz::CaseHandling::STRICTLY_CASE_SENSITIVE,
    morfeusz::CaseHandling::IGNORE_CASE,
};
const morfeusz::WhitespaceHandling translateWhitespaceHandling[] = {
    morfeusz::WhitespaceHandling::SKIP_WHITESPACES,
    morfeusz::WhitespaceHandling::APPEND_WHITESPACES,
    morfeusz::WhitespaceHandling::KEEP_WHITESPACES,
};
const morfeusz::MorfeuszUsage translateUsage[] = {
    morfeusz::MorfeuszUsage::BOTH_ANALYSE_AND_GENERATE,
    morfeusz::MorfeuszUsage::ANALYSE_ONLY,
    morfeusz::MorfeuszUsage::GENERATE_ONLY,
};

const int invalidId = -1;
const struct String emptyString = {};
const Error noError = emptyString;
const struct StringArray emptyStringArray = {};
const struct TokenInfo emptyTokenInfo = {};

template<typename R, typename T, int N>
R reverseTranslate(const T (&array)[N], T value) {
  for (int i = 0; i < N; ++i) {
    if (array[i] == value) {
      return static_cast<R>(i);
    }
  }
  return static_cast<R>(invalidId);
}

const struct String makeString(const char* p, int n) {
  char* cp = new char[n];
  memcpy(cp, p, n);
  return { cp, n };
}

const struct String makeString(const std::string& s) {
  return makeString(s.data(), s.size());
}

const Error makeError(const std::exception& e) {
  return makeString(e.what(), strlen(e.what()));
}

const std::string stdString(const struct String s) {
  // s.p points to the character array underlying a Go string
  // so it will be freed by the Go garbage collector.
  return std::string(s.p, s.n);
}

const struct TokenInfo makeTokenInfo(const MorphInterpretation& m) {
  return {
      makeString(m.orth),
      makeString(m.lemma),
      m.startNode,
      m.endNode,
      m.tagId,
      m.nameId,
      m.labelsId,
  };
}

template<typename T>
const struct StringArray makeStringArray(const T& lst) {
  const int n = lst.size();
  struct String* sp = new struct String[n];
  const struct StringArray ret = { sp, n };
  for (typename T::const_iterator it = lst.begin(); it != lst.end(); ++it) {
    *sp++ = makeString(*it);
  }
  return ret;
}

const struct TokenInfoArray makeTokenInfoArray(
    const std::vector<MorphInterpretation>& vec) {
  const int n = vec.size();
  struct TokenInfo* tp = new struct TokenInfo[n];
  const struct TokenInfoArray ret = { tp, n, noError };
  for (std::vector<MorphInterpretation>::const_iterator it = vec.begin();
       it != vec.end(); ++it) {
    *tp++ = makeTokenInfo(*it);
  }
  return ret;
}

const struct TokenInfoArray makeTokenInfoArray(const std::exception& e) {
  return { NULL, 0, makeError(e) };
}

const Morfeusz* cmcast(const Morf m) {
  return static_cast<const Morfeusz*>(m);
}

Morfeusz* mcast(Morf m) {
  return static_cast<Morfeusz*>(m);
}

ResultsIterator* rcast(Res r) {
  return static_cast<ResultsIterator*>(r);
}

const IdResolver& idResolver(const Morf m) {
  return cmcast(m)->getIdResolver();
}

}  // namespace

extern "C" {

Morf createInstance(const struct String dictName, enum Usage usage) {
  try {
    const morfeusz::MorfeuszUsage morfeuszUsage = translateUsage[usage];
    if (dictName.p == NULL) {
      return Morfeusz::createInstance(morfeuszUsage);
    } else {
      return Morfeusz::createInstance(stdString(dictName), morfeuszUsage);
    }
  } catch (const std::exception& e) {
    return NULL;
  }
}

Res analyseString(const Morf m, const struct String text) {
  try {
    return cmcast(m)->analyse(stdString(text));
  } catch (const std::exception&) {
    return NULL;
  }
}

int hasNext(Res r) {
  return rcast(r)->hasNext();
}

const struct TokenInfo next(Res r) {
  try {
    return makeTokenInfo(rcast(r)->next());
  } catch (const std::exception&) {
    return emptyTokenInfo;
  }
}

const struct String tagsetId(const Morf m) {
  return makeString(idResolver(m).getTagsetId());
}

const struct String tag(const Morf m, int tagId) {
  try {
    return makeString(idResolver(m).getTag(tagId));
  } catch (const std::exception&) {
    return emptyString;
  }
}

int tagId(const Morf m, struct String tag) {
  try {
    return idResolver(m).getTagId(stdString(tag));
  } catch (const std::exception&) {
    return invalidId;
  }
}

const struct String name(const Morf m, int nameId) {
  try {
    return makeString(idResolver(m).getName(nameId));
  } catch (const std::exception&) {
    return emptyString;
  }
}

int nameId(const Morf m, const struct String name) {
  try {
    return idResolver(m).getNameId(stdString(name));
  } catch (const std::exception&) {
    return invalidId;
  }
}

const struct String labelsAsString(const Morf m, int labelsId) {
  try {
    return makeString(idResolver(m).getLabelsAsString(labelsId));
  } catch (const std::exception&) {
    return emptyString;
  }
}

const struct StringArray labels(const Morf m, int labelsId) {
  try {
    return makeStringArray(idResolver(m).getLabels(labelsId));
  } catch (const std::exception&) {
    return emptyStringArray;
  }
}

int labelsId(const Morf m, const struct String labels) {
  try {
    return idResolver(m).getLabelsId(stdString(labels));
  } catch (const std::exception&) {
    return invalidId;
  }
}

int tagsCount(const Morf m) {
  return idResolver(m).getTagsCount();
}

int namesCount(const Morf m) {
  return idResolver(m).getNamesCount();
}

int labelsCount(const Morf m) {
  return idResolver(m).getLabelsCount();
}

const struct TokenInfoArray generate(const Morf m, const struct String lemma) {
  try {
    std::vector<MorphInterpretation> vec;
    cmcast(m)->generate(stdString(lemma), vec);
    return makeTokenInfoArray(vec);
  } catch (const std::exception& e) {
    return makeTokenInfoArray(e);
  }
}

const struct TokenInfoArray generateWithTagID(
    const Morf m, int tagId, const struct String lemma) {
  try {
    std::vector<MorphInterpretation> vec;
    cmcast(m)->generate(stdString(lemma), tagId, vec);
    return makeTokenInfoArray(vec);
  } catch (const std::exception& e) {
    return makeTokenInfoArray(e);
  }
}

const struct String dictId(const Morf m) {
  return makeString(cmcast(m)->getDictID());
}

const struct String dictCopyright(const Morf m) {
  return makeString(cmcast(m)->getDictCopyright());
}

const Error setAggl(Morf m, const struct String aggl) {
  try {
    mcast(m)->setAggl(stdString(aggl));
    return noError;
  } catch (const std::exception& e) {
    return makeError(e);
  }
}

const Error setPraet(Morf m, const struct String praet) {
  try {
    mcast(m)->setPraet(stdString(praet));
    return noError;
  } catch (const std::exception& e) {
    return makeError(e);
  }
}

const Error setCharset(Morf m, enum Charset encoding) {
  try {
    mcast(m)->setCharset(translateCharset[encoding]);
    return noError;
  } catch (const std::exception& e) {
    return makeError(e);
  }
}

const Error setCaseHandling(Morf m, enum CaseHandling caseHandling) {
  try {
    mcast(m)->setCaseHandling(translateCaseHandling[caseHandling]);
    return noError;
  } catch (const std::exception& e) {
    return makeError(e);
  }
}

const Error setTokenNumbering(Morf m, enum TokenNumbering numbering) {
  try {
    mcast(m)->setTokenNumbering(translateTokenNumbering[numbering]);
    return noError;
  } catch (const std::exception& e) {
    return makeError(e);
  }
}

const Error setWhitespaceHandling(Morf m, enum WhitespaceHandling handling) {
  try {
    mcast(m)->setWhitespaceHandling(translateWhitespaceHandling[handling]);
    return noError;
  } catch (const std::exception& e) {
    return makeError(e);
  }
}

const Error setDictionary(Morf m, const struct String dictName) {
  try {
    mcast(m)->setDictionary(stdString(dictName));
    return noError;
  } catch (const std::exception& e) {
    return makeError(e);
  }
}

void setDebug(Morf m, int debug) {
  mcast(m)->setDebug(debug);
}

const struct String aggl(const Morf m) {
  return makeString(cmcast(m)->getAggl());
}

const struct String praet(const Morf m) {
  return makeString(cmcast(m)->getPraet());
}

enum Charset charset(const Morf m) {
  return reverseTranslate<Charset>(
      translateCharset, cmcast(m)->getCharset());
}

enum CaseHandling caseHandling(const Morf m) {
  return reverseTranslate<CaseHandling>(
      translateCaseHandling, cmcast(m)->getCaseHandling());
}

enum TokenNumbering tokenNumbering(const Morf m) {
  return reverseTranslate<TokenNumbering>(
      translateTokenNumbering, cmcast(m)->getTokenNumbering());
}

enum WhitespaceHandling whitespaceHandling(const Morf m) {
  return reverseTranslate<WhitespaceHandling>(
      translateWhitespaceHandling, cmcast(m)->getWhitespaceHandling());
}

const struct StringArray availableAgglOptions(const Morf m) {
  return makeStringArray(cmcast(m)->getAvailableAgglOptions());
}

const struct StringArray availablePraetOptions(const Morf m) {
  return makeStringArray(cmcast(m)->getAvailablePraetOptions());
}

const struct StringArray dictionarySearchPaths(Morf m) {
  return makeStringArray(mcast(m)->dictionarySearchPaths);
}

void prependToDictionarySearchPaths(Morf m, const struct String path) {
  mcast(m)->dictionarySearchPaths.push_front(stdString(path));
}

void appendToDictionarySearchPaths(Morf m, const struct String path) {
  mcast(m)->dictionarySearchPaths.push_back(stdString(path));
}

int removeFromDictionarySearchPaths(Morf m, const struct String path) {
  std::list<std::string>& dsp = mcast(m)->dictionarySearchPaths;
  const size_t previousLength = dsp.size();
  dsp.remove(stdString(path));
  return previousLength - dsp.size();
}

void clearDictionarySearchPaths(Morf m) {
  mcast(m)->dictionarySearchPaths.clear();
}

Morf cloneMorf(const Morf m) {
  return cmcast(m)->clone();
}

void freeMorf(const Morf m) {
  delete cmcast(m);
}

void freeRes(const Res r) {
  delete rcast(r);
}

void freeTokenInfo(const struct TokenInfo* t) {
  delete[] t->orth.p;
  delete[] t->lemma.p;
}

void freeStringArray(const struct StringArray* arr) {
  // The calls to freeCharArray(arr->strings[i].p) happen earlier,
  // when the elements are converted to Go strings via goStringFree().
  delete[] arr->strings;
}

void freeTokenInfoArray(const struct TokenInfoArray* arr) {
  // The calls to freeTokenInfo(arr->tokens[i]) happen later,
  // once the elements become inaccessible.
  delete[] arr->tokens;
  delete[] arr->error.p;
}

void freeCharArray(const char* p) {
  delete[] p;
}

const struct String version() {
  return makeString(Morfeusz::getVersion());
}

const struct String defaultDictName() {
  return makeString(Morfeusz::getDefaultDictName());
}

const struct String copyright() {
  return makeString(Morfeusz::getCopyright());
}

}  // extern "C"
