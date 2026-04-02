package com.lieying.platform.controller;

import com.lieying.platform.model.Task;
import com.lieying.platform.repository.TaskRepository;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.tags.Tag;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.UUID;

@RestController
@RequestMapping("/api/v1/tasks")
@RequiredArgsConstructor
@Tag(name = "Task Management", description = "任务管理API")
public class TaskController {
    private final TaskRepository taskRepository;

    @PostMapping
    @Operation(summary = "创建任务")
    public ResponseEntity<Task> create(@RequestBody Task task) {
        return ResponseEntity.ok(taskRepository.save(task));
    }

    @GetMapping
    @Operation(summary = "获取任务列表")
    public ResponseEntity<List<Task>> list() {
        return ResponseEntity.ok(taskRepository.findAll());
    }

    @GetMapping("/{id}")
    @Operation(summary = "获取任务详情")
    public ResponseEntity<Task> get(@PathVariable UUID id) {
        return taskRepository.findById(id)
                .map(ResponseEntity::ok)
                .orElse(ResponseEntity.notFound().build());
    }
}
