package com.lieying.platform.controller;

import com.lieying.platform.service.AiService;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.tags.Tag;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/api/v1/ai")
@RequiredArgsConstructor
@Tag(name = "AI Service", description = "AI服务API")
public class AiController {
    private final AiService aiService;

    @PostMapping("/poc/generate")
    @Operation(summary = "AI生成POC")
    public ResponseEntity<String> generatePoc(@RequestBody String vulnDescription) {
        return ResponseEntity.ok(aiService.generatePoc(vulnDescription));
    }

    @PostMapping("/report/generate")
    @Operation(summary = "AI生成报告")
    public ResponseEntity<String> generateReport(@RequestBody String vulnDetails) {
        return ResponseEntity.ok(aiService.generateReport(vulnDetails));
    }

    @PostMapping("/attack-path/analyze")
    @Operation(summary = "AI分析攻击路径")
    public ResponseEntity<String> analyzeAttackPath(@RequestBody String reconData) {
        return ResponseEntity.ok(aiService.analyzeAttackPath(reconData));
    }
}
