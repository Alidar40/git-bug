// Code generated by "stringer -type TokenType"; DO NOT EDIT.

package chroma

import "strconv"

const _TokenType_name = "NoneOtherErrorLineTableTDLineTableLineHighlightLineNumbersTableLineNumbersBackgroundEOFTypeKeywordKeywordConstantKeywordDeclarationKeywordNamespaceKeywordPseudoKeywordReservedKeywordTypeNameNameAttributeNameBuiltinNameBuiltinPseudoNameClassNameConstantNameDecoratorNameEntityNameExceptionNameFunctionNameFunctionMagicNameKeywordNameLabelNameNamespaceNameOperatorNameOtherNamePseudoNamePropertyNameTagNameVariableNameVariableAnonymousNameVariableClassNameVariableGlobalNameVariableInstanceNameVariableMagicLiteralLiteralDateLiteralOtherLiteralStringLiteralStringAffixLiteralStringAtomLiteralStringBacktickLiteralStringBooleanLiteralStringCharLiteralStringDelimiterLiteralStringDocLiteralStringDoubleLiteralStringEscapeLiteralStringHeredocLiteralStringInterpolLiteralStringNameLiteralStringOtherLiteralStringRegexLiteralStringSingleLiteralStringSymbolLiteralNumberLiteralNumberBinLiteralNumberFloatLiteralNumberHexLiteralNumberIntegerLiteralNumberIntegerLongLiteralNumberOctOperatorOperatorWordPunctuationCommentCommentHashbangCommentMultilineCommentSingleCommentSpecialCommentPreprocCommentPreprocFileGenericGenericDeletedGenericEmphGenericErrorGenericHeadingGenericInsertedGenericOutputGenericPromptGenericStrongGenericSubheadingGenericTracebackGenericUnderlineTextTextWhitespaceTextSymbolTextPunctuation"

var _TokenType_map = map[TokenType]string{
	-9:   _TokenType_name[0:4],
	-8:   _TokenType_name[4:9],
	-7:   _TokenType_name[9:14],
	-6:   _TokenType_name[14:25],
	-5:   _TokenType_name[25:34],
	-4:   _TokenType_name[34:47],
	-3:   _TokenType_name[47:63],
	-2:   _TokenType_name[63:74],
	-1:   _TokenType_name[74:84],
	0:    _TokenType_name[84:91],
	1000: _TokenType_name[91:98],
	1001: _TokenType_name[98:113],
	1002: _TokenType_name[113:131],
	1003: _TokenType_name[131:147],
	1004: _TokenType_name[147:160],
	1005: _TokenType_name[160:175],
	1006: _TokenType_name[175:186],
	2000: _TokenType_name[186:190],
	2001: _TokenType_name[190:203],
	2002: _TokenType_name[203:214],
	2003: _TokenType_name[214:231],
	2004: _TokenType_name[231:240],
	2005: _TokenType_name[240:252],
	2006: _TokenType_name[252:265],
	2007: _TokenType_name[265:275],
	2008: _TokenType_name[275:288],
	2009: _TokenType_name[288:300],
	2010: _TokenType_name[300:317],
	2011: _TokenType_name[317:328],
	2012: _TokenType_name[328:337],
	2013: _TokenType_name[337:350],
	2014: _TokenType_name[350:362],
	2015: _TokenType_name[362:371],
	2016: _TokenType_name[371:381],
	2017: _TokenType_name[381:393],
	2018: _TokenType_name[393:400],
	2019: _TokenType_name[400:412],
	2020: _TokenType_name[412:433],
	2021: _TokenType_name[433:450],
	2022: _TokenType_name[450:468],
	2023: _TokenType_name[468:488],
	2024: _TokenType_name[488:505],
	3000: _TokenType_name[505:512],
	3001: _TokenType_name[512:523],
	3002: _TokenType_name[523:535],
	3100: _TokenType_name[535:548],
	3101: _TokenType_name[548:566],
	3102: _TokenType_name[566:583],
	3103: _TokenType_name[583:604],
	3104: _TokenType_name[604:624],
	3105: _TokenType_name[624:641],
	3106: _TokenType_name[641:663],
	3107: _TokenType_name[663:679],
	3108: _TokenType_name[679:698],
	3109: _TokenType_name[698:717],
	3110: _TokenType_name[717:737],
	3111: _TokenType_name[737:758],
	3112: _TokenType_name[758:775],
	3113: _TokenType_name[775:793],
	3114: _TokenType_name[793:811],
	3115: _TokenType_name[811:830],
	3116: _TokenType_name[830:849],
	3200: _TokenType_name[849:862],
	3201: _TokenType_name[862:878],
	3202: _TokenType_name[878:896],
	3203: _TokenType_name[896:912],
	3204: _TokenType_name[912:932],
	3205: _TokenType_name[932:956],
	3206: _TokenType_name[956:972],
	4000: _TokenType_name[972:980],
	4001: _TokenType_name[980:992],
	5000: _TokenType_name[992:1003],
	6000: _TokenType_name[1003:1010],
	6001: _TokenType_name[1010:1025],
	6002: _TokenType_name[1025:1041],
	6003: _TokenType_name[1041:1054],
	6004: _TokenType_name[1054:1068],
	6100: _TokenType_name[1068:1082],
	6101: _TokenType_name[1082:1100],
	7000: _TokenType_name[1100:1107],
	7001: _TokenType_name[1107:1121],
	7002: _TokenType_name[1121:1132],
	7003: _TokenType_name[1132:1144],
	7004: _TokenType_name[1144:1158],
	7005: _TokenType_name[1158:1173],
	7006: _TokenType_name[1173:1186],
	7007: _TokenType_name[1186:1199],
	7008: _TokenType_name[1199:1212],
	7009: _TokenType_name[1212:1229],
	7010: _TokenType_name[1229:1245],
	7011: _TokenType_name[1245:1261],
	8000: _TokenType_name[1261:1265],
	8001: _TokenType_name[1265:1279],
	8002: _TokenType_name[1279:1289],
	8003: _TokenType_name[1289:1304],
}

func (i TokenType) String() string {
	if str, ok := _TokenType_map[i]; ok {
		return str
	}
	return "TokenType(" + strconv.FormatInt(int64(i), 10) + ")"
}